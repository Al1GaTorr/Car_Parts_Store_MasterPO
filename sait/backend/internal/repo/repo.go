package repo

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"bazarpo-backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	DB     *mongo.Database
	Users  *mongo.Collection
	Cars   *mongo.Collection
	Parts  *mongo.Collection
	Orders *mongo.Collection
}

func New(db *mongo.Database) *Repository {
	return &Repository{
		DB:     db,
		Users:  db.Collection("users"),
		Cars:   db.Collection("cars"),
		Parts:  db.Collection("parts"),
		Orders: db.Collection("orders"),
	}
}

func (r *Repository) CountUsersByEmail(ctx context.Context, email string) (int64, error) {
	return r.Users.CountDocuments(ctx, bson.M{"email": email})
}

func (r *Repository) InsertUser(ctx context.Context, user model.UserDoc) (primitive.ObjectID, error) {
	res, err := r.Users.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*model.UserDoc, error) {
	var user model.UserDoc
	if err := r.Users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*model.UserDoc, error) {
	var user model.UserDoc
	if err := r.Users.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUserRoleByID(ctx context.Context, id primitive.ObjectID, role string) error {
	_, err := r.Users.UpdateByID(ctx, id, bson.M{"$set": bson.M{"role": role}})
	return err
}

func (r *Repository) FindCarByVIN(ctx context.Context, vin string) (*model.CarDoc, error) {
	var car model.CarDoc
	if err := r.Cars.FindOne(ctx, bson.M{"vin": vin}).Decode(&car); err != nil {
		return nil, err
	}
	return &car, nil
}

func exactCI(value string) bson.M {
	return bson.M{
		"$regex":   "^" + regexp.QuoteMeta(strings.TrimSpace(value)) + "$",
		"$options": "i",
	}
}

func distinctStrings(values []any) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		s, ok := v.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func distinctInts(values []any) ([]int, error) {
	out := make([]int, 0, len(values))
	for _, v := range values {
		switch n := v.(type) {
		case int:
			out = append(out, n)
		case int32:
			out = append(out, int(n))
		case int64:
			out = append(out, int(n))
		case float64:
			out = append(out, int(n))
		default:
			return nil, fmt.Errorf("unsupported year type %T", v)
		}
	}
	sort.Ints(out)
	return out, nil
}

func (r *Repository) DistinctCarMakes(ctx context.Context) ([]string, error) {
	values, err := r.Cars.Distinct(ctx, "make", bson.M{})
	if err != nil {
		return nil, err
	}
	return distinctStrings(values), nil
}

func (r *Repository) DistinctCarModelsByMake(ctx context.Context, make string) ([]string, error) {
	values, err := r.Cars.Distinct(ctx, "model", bson.M{"make": exactCI(make)})
	if err != nil {
		return nil, err
	}
	return distinctStrings(values), nil
}

func (r *Repository) DistinctCarYearsByMakeModel(ctx context.Context, make, model string) ([]int, error) {
	filter := bson.M{
		"make":  exactCI(make),
		"model": exactCI(model),
	}
	values, err := r.Cars.Distinct(ctx, "year", filter)
	if err != nil {
		return nil, err
	}
	return distinctInts(values)
}

func (r *Repository) FindParts(ctx context.Context, filter bson.M, limit int64) ([]model.PartDoc, error) {
	cur, err := r.Parts.Find(ctx, filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []model.PartDoc
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) FindPartsBySKUs(ctx context.Context, skus []string) ([]model.PartDoc, error) {
	if len(skus) == 0 {
		return []model.PartDoc{}, nil
	}
	cur, err := r.Parts.Find(ctx, bson.M{"sku": bson.M{"$in": skus}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []model.PartDoc
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) DecrementPartStockIfEnough(ctx context.Context, sku string, qty int) (bool, error) {
	res, err := r.Parts.UpdateOne(
		ctx,
		bson.M{
			"sku":       strings.TrimSpace(sku),
			"stock_qty": bson.M{"$gte": qty},
		},
		bson.M{
			"$inc": bson.M{"stock_qty": -qty},
		},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *Repository) IncrementPartStock(ctx context.Context, sku string, qty int) error {
	_, err := r.Parts.UpdateOne(
		ctx,
		bson.M{"sku": strings.TrimSpace(sku)},
		bson.M{"$inc": bson.M{"stock_qty": qty}},
	)
	return err
}

func (r *Repository) FindPartBySKU(ctx context.Context, sku string) (*model.PartDoc, error) {
	var p model.PartDoc
	if err := r.Parts.FindOne(ctx, bson.M{"sku": strings.TrimSpace(sku)}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListOrders(ctx context.Context, limit int64) ([]model.OrderDoc, error) {
	cur, err := r.Orders.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []model.OrderDoc
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) InsertOrder(ctx context.Context, order model.OrderDoc) (primitive.ObjectID, error) {
	res, err := r.Orders.InsertOne(ctx, order)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func (r *Repository) UpdateOrderByID(ctx context.Context, id primitive.ObjectID, set bson.M) error {
	_, err := r.Orders.UpdateByID(ctx, id, bson.M{"$set": set})
	return err
}

func (r *Repository) DeleteOrderByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Orders.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *Repository) ListPartsForAdmin(ctx context.Context, limit int64) ([]model.PartDoc, error) {
	cur, err := r.Parts.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []model.PartDoc
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) InsertPart(ctx context.Context, p model.PartDoc) error {
	_, err := r.Parts.InsertOne(ctx, p)
	return err
}

func (r *Repository) UpdatePartByID(ctx context.Context, id primitive.ObjectID, set bson.M) error {
	_, err := r.Parts.UpdateByID(ctx, id, bson.M{"$set": set})
	return err
}

func (r *Repository) DeletePartByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Parts.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
