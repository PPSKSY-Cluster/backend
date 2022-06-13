package db

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CResourceType int64

const (
	customRT CResourceType = -1
	noneRT   CResourceType = 0
	gpuRT    CResourceType = 1
)

type LoadBalancing int64

const (
	customLB     LoadBalancing = -1
	noneLB       LoadBalancing = 0
	roundRobinLB LoadBalancing = 1
)

type OSType int64

const (
	linuxOS   OSType = 1
	windowsOS OSType = 2
	macOS     OSType = 3
)

type CResource struct {
	ID                       primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                     string               `bson:"name" json:"name"`
	Description              string               `bson:"description" json:"description"`
	Nodes                    []int64              `bson:"nodes" json:"nodes"`
	Type                     CResourceType        `bson:"type" json:"type"`
	Admins                   []User               `bson:"admins" json:"admins"`
	BalancingAlg             LoadBalancing        `bson:"balancingAlg" json:"balancingAlg"`
	Reservations             []primitive.ObjectID `bson:"reservations" json:"reservations"`
	HighAvailability         bool                 `bson:"highAvailability" json:"highAvailability"`                 // HA
	HighPerformanceComputing bool                 `bson:"highPerformanceComputing" json:"highPerformanceComputing"` // HPC
	OperatingSystem          OSType               `bson:"operatingSystem" json:"operatingSystem"`
}

var collectionName = "cResources"

func GetAllResources() ([]CResource, error) {

	query := func() (*mongo.Cursor, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection(collectionName).
			Find(mdbInstance.Ctx, bson.M{})
	}

	resourceCursor, err := runQuery[*mongo.Cursor](query)
	if err != nil {
		return nil, err
	}

	var cResources []CResource
	if err = resourceCursor.All(mdbInstance.Ctx, &cResources); err != nil {
		return nil, err
	}

	return cResources, nil
}

func GetResourceById(_id primitive.ObjectID) (CResource, error) {

	query := func() (*mongo.SingleResult, error) {
		singleRes := mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection(collectionName).
			FindOne(mdbInstance.Ctx, bson.M{"_id": _id})
		return singleRes, singleRes.Err()
	}

	resourceSingleRes, err := runQuery[*mongo.SingleResult](query)
	if err != nil {
		return CResource{}, err
	}

	var cResource CResource
	if err = resourceSingleRes.Decode(&cResource); err != nil {
		return CResource{}, err
	}

	return cResource, nil
}

func AddResource(cResource CResource) (CResource, error) {

	query := func() (*mongo.InsertOneResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection(collectionName).
			InsertOne(mdbInstance.Ctx, cResource)
	}

	if err := mdbInstance.Validate.Struct(cResource); err != nil {
		return CResource{}, err
	}

	insertRes, err := runQuery[*mongo.InsertOneResult](query)
	if err != nil {
		return CResource{}, err
	}

	cResource.ID = insertRes.InsertedID.(primitive.ObjectID)

	return cResource, nil
}

func EditResource(_id primitive.ObjectID, cResource CResource) (CResource, error) {

	query := func() (*mongo.UpdateResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection(collectionName).
			UpdateOne(mdbInstance.Ctx,
				bson.M{"_id": _id},
				bson.M{"$set": cResource})
	}

	_, err := runQuery[*mongo.UpdateResult](query)
	if err != nil {
		return CResource{}, err
	}

	return cResource, nil
}

func DeleteResource(_id primitive.ObjectID) error {

	query := func() (*mongo.DeleteResult, error) {
		return mdbInstance.Client.
			Database(os.Getenv("DB_NAME")).
			Collection(collectionName).
			DeleteOne(mdbInstance.Ctx, bson.M{"_id": _id})
	}

	_, err := runQuery[*mongo.DeleteResult](query)
	if err != nil {
		return err
	}

	return nil
}
