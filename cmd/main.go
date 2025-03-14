package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	batch = 100
)

func copyData(dbSrc *mongo.Database, dbDst *mongo.Database, collection string) {
	counter := 0

	collSrc := dbSrc.Collection(collection)
	collDst := dbDst.Collection(collection)

	total, _ := collSrc.EstimatedDocumentCount(context.TODO())
	log.Printf("Copying %d documents for collection %s", total, collection)

	documents, err := collSrc.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Cannot retrieve documents", err)
		return
	}

	insert := make([]bson.Raw, batch)
	opts := options.InsertMany().SetOrdered(false)

	for documents.Next(context.TODO()) {
		insert[counter%batch] = documents.Current

		counter++
		if counter%batch == 0 {
			result, err := collDst.InsertMany(context.TODO(), insert, opts)

			if result == nil && err != nil {
				log.Println("Cannot copy documents to collection", collection, err)
				return
			}

			log.Printf("Copied %d / %d documents to collection %s, progress: %d / %d", len(result.InsertedIDs), batch, collection, counter, total)
		}
	}

	if counter%batch != 0 {
		collDst.InsertMany(context.TODO(), insert[0:counter%batch], opts)
	}

	log.Printf("Finished copying %d documents from collection %s", counter, collection)

	documents.Close(context.TODO())
}

func main() {
	srcUri := flag.String("src", "", "Specify MongoDB URL source")
	dstUri := flag.String("dst", "", "Specify MongoDB URL destination")
	flag.IntVar(&batch, "batch", 100, "Insert Batch size")

	flag.Parse()

	if *srcUri == "" || *dstUri == "" {
		log.Fatal("Please specify source and destination URL")
	}

	srcOpts := options.Client().ApplyURI(*srcUri)
	srcClient, err := mongo.Connect(srcOpts)
	if err != nil {
		log.Panic(err)
	}

	dstOpts := options.Client().ApplyURI(*dstUri)
	dstClient, err := mongo.Connect(dstOpts)
	if err != nil {
		log.Panic(err)
	}

	parsedURL, err := url.Parse(srcOpts.GetURI())
	if err != nil {
		log.Fatal(err)
	}
	dbSrc := srcClient.Database(parsedURL.Path[1:])

	parsedURL, err = url.Parse(dstOpts.GetURI())
	if err != nil {
		log.Fatal(err)
	}
	dbDst := dstClient.Database(parsedURL.Path[1:])

	collections, err := dbSrc.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Collections ", collections)

	var wg sync.WaitGroup

	for _, col := range collections {
		wg.Add(1)
		go func() {
			defer wg.Done()
			copyData(dbSrc, dbDst, col)
		}()
	}

	wg.Wait()

	dstClient.Disconnect(context.TODO())
	srcClient.Disconnect(context.TODO())
}
