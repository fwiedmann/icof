package gocrazy_permanent_data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/icof/pkg/notifier"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"io"
	"net/http"
)

var bucketName = []byte("email-config")
var configKey = []byte("config")
var dbLocation = "/perm/icof/email-config.db"

// NewDynamicEmailConfig constructs a DynamicEmailConfig instance
func NewDynamicEmailConfig(logger *logrus.Logger) (*DynamicEmailConfig, error) {
	db, err := bbolt.Open(dbLocation, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &DynamicEmailConfig{db: db, logger: logger}, nil
}

// DynamicEmailConfig struct implements the notifier.EmailReceiverRepository
type DynamicEmailConfig struct {
	db     *bbolt.DB
	logger *logrus.Logger
}

// GetEmailReceivers lookups the current config in the bolt database
func (d *DynamicEmailConfig) GetEmailReceivers(ctx context.Context) ([]notifier.EmailReceiver, error) {
	var receivers []notifier.EmailReceiver

	err := d.db.View(func(tx *bbolt.Tx) error {
		config := tx.Bucket(bucketName).Get(configKey)
		if config == nil {
			return fmt.Errorf("could not load config from database")
		}

		return json.Unmarshal(config, &receivers)
	})

	if err != nil {
		return nil, err
	}
	return receivers, nil
}

// CreatOrUpdateConfig validate the input and replaces the current config in the databse
func (d *DynamicEmailConfig) CreatOrUpdateConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			d.logger.Errorf("could read body: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		var receivers []notifier.EmailReceiver
		err = json.Unmarshal(body, &receivers)
		if err != nil {
			d.logger.Errorf("could not unmarshal body: %s", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		v := validator.New()
		for _, r := range receivers {
			err := v.Struct(r)
			if err != nil {
				d.logger.Errorf("invalid email receiver: %s", err.Error())
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}

		err = d.db.Update(func(tx *bbolt.Tx) error {
			return tx.Bucket(bucketName).Put(configKey, body)
		})

		if err != nil {
			d.logger.Errorf("could not update bucket item: %s", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// GetConfigHandler lookups the current config in the DB and writes it as json to the client
func (d *DynamicEmailConfig) GetConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config []byte
		err := d.db.View(func(tx *bbolt.Tx) error {
			config = tx.Bucket(bucketName).Get(configKey)
			if config == nil {
				d.logger.Errorf("could not find config in database")
				return fmt.Errorf("could not load config from database")
			}
			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(config)
		if err != nil {
			d.logger.Errorf("error while writing response to client: %s", err.Error())
		}
	}
}
