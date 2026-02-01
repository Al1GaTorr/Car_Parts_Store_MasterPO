package main

import (
	"context"
	"errors"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	categories *mongo.Collection
	parts      *mongo.Collection
	orders     *mongo.Collection
	alerts     *mongo.Collection

	lowStockCh chan LowStockAlert
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		categories: db.Collection("categories"),
		parts:      db.Collection("spare_parts"),
		orders:     db.Collection("orders"),
		alerts:     db.Collection("alerts"),
		lowStockCh: make(chan LowStockAlert, 100),
	}
}

// -------- categories --------
func (r *Repo) CreateCategory(ctx context.Context, c Category) (Category, error) {
	res, err := r.categories.InsertOne(ctx, c)
	if err != nil {
		return Category{}, err
	}
	c.ID = res.InsertedID.(primitive.ObjectID)
	return c, nil
}

func (r *Repo) ListCategories(ctx context.Context) ([]Category, error) {
	cur, err := r.categories.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]Category, 0)
	for cur.Next(ctx) {
		var c Category
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *Repo) GetCategory(ctx context.Context, id primitive.ObjectID) (Category, error) {
	var c Category
	err := r.categories.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	return c, err
}

func (r *Repo) UpdateCategory(ctx context.Context, id primitive.ObjectID, name, desc string) (Category, error) {
	upd := bson.M{}
	if name != "" {
		upd["name"] = name
	}
	if desc != "" {
		upd["description"] = desc
	}
	if len(upd) == 0 {
		return Category{}, errors.New("nothing to update")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out Category
	err := r.categories.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": upd}, opts).Decode(&out)
	return out, err
}

func (r *Repo) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.categories.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// -------- parts --------
func (r *Repo) CreatePart(ctx context.Context, p SparePart) (SparePart, error) {
	res, err := r.parts.InsertOne(ctx, p)
	if err != nil {
		return SparePart{}, err
	}
	p.ID = res.InsertedID.(primitive.ObjectID)

	// optional: push part to category parts_list
	_, _ = r.categories.UpdateOne(ctx, bson.M{"_id": p.CategoryID}, bson.M{"$addToSet": bson.M{"parts_list": p.ID}})
	return p, nil
}

func (r *Repo) GetPart(ctx context.Context, id primitive.ObjectID) (SparePart, error) {
	var p SparePart
	err := r.parts.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	return p, err
}

func (r *Repo) DeletePart(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.parts.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repo) UpdatePart(ctx context.Context, id primitive.ObjectID, upd bson.M) (SparePart, error) {
	if len(upd) == 0 {
		return SparePart{}, errors.New("nothing to update")
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out SparePart
	err := r.parts.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": upd}, opts).Decode(&out)
	return out, err
}

func (r *Repo) ListPartsFiltered(ctx context.Context, categoryID *primitive.ObjectID, carModel, brand, q, compatibility string) ([]SparePart, error) {
	filter := bson.M{"is_active": true}

	if categoryID != nil {
		filter["category_id"] = *categoryID
	}
	if carModel != "" {
		filter["car_model"] = bson.M{"$regex": regexp.QuoteMeta(carModel), "$options": "i"}
	}
	if brand != "" {
		filter["brand"] = bson.M{"$regex": regexp.QuoteMeta(brand), "$options": "i"}
	}
	if compatibility != "" {
		filter["compatibility"] = bson.M{"$regex": regexp.QuoteMeta(compatibility), "$options": "i"}
	}
	if q != "" {
		re := bson.M{"$regex": regexp.QuoteMeta(q), "$options": "i"}
		filter["$or"] = []bson.M{
			{"description": re},
			{"brand": re},
			{"car_model": re},
		}
	}

	cur, err := r.parts.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]SparePart, 0)
	for cur.Next(ctx) {
		var p SparePart
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// DecreaseStock: atomic check + decrement
func (r *Repo) DecreaseStock(ctx context.Context, partID primitive.ObjectID, qty int) (SparePart, error) {
	if qty <= 0 {
		return SparePart{}, errors.New("quantity must be > 0")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated SparePart

	filter := bson.M{
		"_id":       partID,
		"is_active": true,
		"stock":     bson.M{"$gte": qty},
	}

	err := r.parts.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$inc": bson.M{"stock": -qty}},
		opts,
	).Decode(&updated)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return SparePart{}, errors.New("not enough stock or part not found")
		}
		return SparePart{}, err
	}

	if updated.Stock <= 5 {
		r.lowStockCh <- LowStockAlert{
			PartID: updated.ID,
			Name:   updated.Brand + " " + updated.CarModel,
			Stock:  updated.Stock,
			At:     time.Now(),
		}
	}

	return updated, nil
}

// -------- orders --------
func (r *Repo) CreateOrder(ctx context.Context, o Order) (Order, error) {
	res, err := r.orders.InsertOne(ctx, o)
	if err != nil {
		return Order{}, err
	}
	o.ID = res.InsertedID.(primitive.ObjectID)
	return o, nil
}

func (r *Repo) GetOrder(ctx context.Context, id primitive.ObjectID) (Order, error) {
	var o Order
	err := r.orders.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	return o, err
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, id primitive.ObjectID, status string, isPaid bool) (Order, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out Order
	err := r.orders.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status, "is_paid": isPaid}}, opts).Decode(&out)
	return out, err
}

func (r *Repo) CancelOrder(ctx context.Context, id primitive.ObjectID) (Order, error) {
	return r.UpdateOrderStatus(ctx, id, "canceled", false)
}

// -------- alerts --------
func (r *Repo) InsertAlert(ctx context.Context, a LowStockAlert) error {
	_, err := r.alerts.InsertOne(ctx, a)
	return err
}

func (r *Repo) ListAlerts(ctx context.Context, limit int64) ([]LowStockAlert, error) {
	opts := options.Find().SetSort(bson.M{"at": -1})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cur, err := r.alerts.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]LowStockAlert, 0)
	for cur.Next(ctx) {
		var a LowStockAlert
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}
