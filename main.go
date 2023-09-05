package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

type Tutorial struct {
	Id       string
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []string
}

type Comment struct {
	Body string
}

func populate() []Tutorial {
	author := &Author{Name: "John Doe", Tutorials: []string{
		"id_1",
	}}

	tutorial := Tutorial{
		Id:     "id_1",
		Title:  "How to make a bread",
		Author: *author,
		Comments: []Comment{
			{Body: "This is a great tutorial"},
			{Body: "Wow, fantastic!"},
			{Body: "Good tutorial"},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials
}

func main() {
	tutorials := populate()

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"Body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
	)

	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"Id": &graphql.Field{
					Type: graphql.String,
				},
				"Title": &graphql.Field{
					Type: graphql.String,
				},
				"Author": &graphql.Field{
					Type: authorType,
				},
				"Comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get tutorial by Id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(string)
				if ok {
					for _, tutorial := range tutorials {
						if tutorial.Id == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Full Tutorial List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	//Defines the object config
	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	// Defines a schema config
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	// Creates our schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL schema, err %v", err.Error())
	}

	// This query returns the "hello" field in the Fields section
	query := `
		{
		  tutorial(id: "id_1") {
			Title
			Comments {
			  Body
			}
			Author {
			  Name
			  Tutorials
			}
		  }
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors %v", err.Error())
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
