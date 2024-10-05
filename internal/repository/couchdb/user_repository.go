package couchdb

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Sokcheatsrorng/go-clean-architecture/internal/model"
	"github.com/go-kivik/kivik/v3"
	
)


type UserRepository struct {
	db *kivik.DB
}

func NewUserRepository(db *kivik.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserCreatePayload struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (r *UserRepository) Create(users []model.User) error {
	docs := make([]interface{}, len(users))
	for i, user := range users {
		docs[i] = UserCreatePayload{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		}
	}

	results, err := r.db.BulkDocs(context.TODO(), docs)
	if err != nil {
		return fmt.Errorf("failed to perform bulk insert: %w", err)
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("failed to close bulk results: %w", err)
	}

	return nil
}

func (r *UserRepository) Read(id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	var user model.User
	err := r.db.Get(context.TODO(), id).ScanDoc(&user)
	if err != nil {
		if kivik.StatusCode(err) == http.StatusNotFound {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving user with id %s: %w", id, err)
	}

	return &user, nil
}

func (r *UserRepository) Update(users []model.User) error {
    if len(users) == 0 {
        return nil
    }

    // Prepare the bulk update documents
    docs := make([]interface{}, len(users))
    for i, user := range users {
        // Validate that each user has an ID and Rev
        if user.ID == "" {
            return fmt.Errorf("user at index %d is missing ID", i)
        }
        if user.Rev == "" {
            return fmt.Errorf("user at index %d is missing Rev", i)
        }

        // Prepare the document to be updated
        doc := map[string]interface{}{
            "_id":      user.ID,
            "_rev":     user.Rev,
            "username": user.Username,
            "email":    user.Email,
            "password": user.Password,
        }

        // Add to bulk update list
        docs[i] = doc
    }

    // Perform the bulk update operation
    results, err := r.db.BulkDocs(context.TODO(), docs)
    if err != nil {
        return fmt.Errorf("failed to perform bulk update: %w", err)
    }
    defer results.Close()

    // Collect and log any update errors
    var updateErrors []error
    for results.Next() {
        if results.Err() != nil {
            updateErrors = append(updateErrors, fmt.Errorf("error updating document %s: %w", results.ID(), results.Err()))
        } else {
            log.Printf("Successfully updated document: %s", results.ID())
        }
    }

    // Return aggregated errors if any occurred during the update process
    if len(updateErrors) > 0 {
        for _, err := range updateErrors {
            log.Printf("Update error: %v", err)
        }
        return fmt.Errorf("encountered %d errors during bulk update", len(updateErrors))
    }

    return nil
}

// Delete removes a user by their ID from the database.
func (r *UserRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	ctx := context.TODO()

	row := r.db.Get(ctx, id)
	if row.Err != nil {
		if kivik.StatusCode(row.Err) == http.StatusNotFound {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("error retrieving user with id %s: %w", id, row.Err)
	}

	var doc struct {
		Rev string `json:"_rev"`
	}
	if err := row.ScanDoc(&doc); err != nil {
		return fmt.Errorf("error scanning document: %w", err)
	}

	_, err := r.db.Delete(ctx, id, doc.Rev)
	if err != nil {
		return fmt.Errorf("error deleting user with id %s: %w", id, err)
	}

	return nil
}

// ReadAll retrieves all users from the database.
func (r *UserRepository) ReadAll() ([]model.User, error) {
    rows, err := r.db.AllDocs(context.TODO(), kivik.Options{"include_docs": true})
    if err != nil {
        return nil, fmt.Errorf("error retrieving all documents: %w", err)
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var doc struct {
            ID  string     `json:"_id"`
            Rev string     `json:"_rev"`
            Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Doc model.User `json:"doc"`
        }
        if err := rows.ScanDoc(&doc); err != nil {
            return nil, fmt.Errorf("failed to scan document: %w", err)
        }
        doc.Doc.ID = doc.ID
        doc.Doc.Rev = doc.Rev
		doc.Doc.Username = doc.Username
		doc.Doc.Email = doc.Email
		doc.Doc.Password = doc.Password
        users = append(users, doc.Doc)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during document iteration: %w", err)
    }

    return users, nil
}