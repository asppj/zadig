/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package repo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/koderover/zadig/pkg/microservice/aslan/config"
	"github.com/koderover/zadig/pkg/microservice/aslan/core/common/dao/models"
	mongotool "github.com/koderover/zadig/pkg/tool/mongo"
)

// BuildPipeResp ...
type BuildPipeResp struct {
	ID           BuildItem `bson:"_id"                    json:"_id"`
	TotalSuccess int       `bson:"total_success"          json:"total_success"`
	TotalFailure int       `bson:"total_failure"          json:"total_failure"`
}

// BuildItem ...
type BuildItem struct {
	ProductName  string `bson:"product_name"       json:"product_name"`
	TotalSuccess int    `bson:"total_success"      json:"total_success"`
	TotalFailure int    `bson:"total_failure"      json:"total_failure"`
}

type BuildStatOption struct {
	StartDate    int64
	EndDate      int64
	Limit        int
	Skip         int
	IsAsc        bool
	ProductNames []string
}

type BuildStatColl struct {
	*mongo.Collection

	coll string
}

func NewBuildStatColl() *BuildStatColl {
	name := models.BuildStat{}.TableName()
	return &BuildStatColl{Collection: mongotool.Database(config.MongoDatabase()).Collection(name), coll: name}
}

func (c *BuildStatColl) GetCollectionName() string {
	return c.coll
}

func (c *BuildStatColl) EnsureIndex(ctx context.Context) error {
	mod := mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "product_name", Value: 1},
			bson.E{Key: "date", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := c.Indexes().CreateOne(ctx, mod)

	return err
}

func (c *BuildStatColl) Create(args *models.BuildStat) error {
	if args == nil {
		return errors.New("nil buildStat args")
	}

	_, err := c.InsertOne(context.TODO(), args)
	return err
}

func (c *BuildStatColl) Update(args *models.BuildStat) error {
	if args == nil {
		return errors.New("nil buildStat args")
	}

	query := bson.M{"date": args.Date, "product_name": args.ProductName}
	update := bson.M{"$set": args}
	_, err := c.UpdateOne(context.TODO(), query, update)
	return err
}

func (c *BuildStatColl) FindCount() (int, error) {
	count, err := c.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (c *BuildStatColl) ListBuildStat(option *BuildStatOption) ([]*models.BuildStat, error) {
	resp := make([]*models.BuildStat, 0)
	query := bson.M{}

	if len(option.ProductNames) > 0 {
		query["product_name"] = bson.M{"$in": option.ProductNames}
	}

	if option.StartDate > 0 {
		query["create_time"] = bson.M{"$gte": option.StartDate, "$lte": option.EndDate}
	}

	opt := &options.FindOptions{}

	if option.Limit > 0 {
		opt.SetSort(bson.D{{"max_duration", -1}}).SetSkip(int64(option.Skip)).SetLimit(int64(option.Limit))
	} else if option.IsAsc {
		opt.SetSort(bson.D{{"create_time", 1}})
	} else {
		opt.SetSort(bson.D{{"create_time", -1}})
	}

	cursor, err := c.Collection.Find(context.TODO(), query, opt)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *BuildStatColl) GetBuildTotalAndSuccess() ([]*BuildItem, error) {
	var result []*BuildPipeResp
	var pipeline []bson.M
	var resp []*BuildItem

	pipeline = append(pipeline,
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"product_name": "$product_name",
				},
				"total_success": bson.M{
					"$sum": "$total_success",
				},
				"total_failure": bson.M{
					"$sum": "$total_failure",
				},
			},
		})

	cursor, err := c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	for _, res := range result {
		buildItem := &BuildItem{
			ProductName:  res.ID.ProductName,
			TotalSuccess: res.TotalSuccess,
			TotalFailure: res.TotalFailure,
		}
		resp = append(resp, buildItem)
	}

	return resp, nil
}
