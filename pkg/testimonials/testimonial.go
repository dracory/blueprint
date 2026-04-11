package testimonials

import (
	"context"
	"errors"

	"github.com/dracory/dataobject"
	"github.com/dracory/entitystore"
)

const ENTITY_TYPE = "testimonial"

// Field constants for testimonial attributes
const (
	FIELD_DATE       = "date"
	FIELD_FIRST_NAME = "first_name"
	FIELD_ID         = "id"
	FIELD_IMAGE_URL  = "image_url"
	FIELD_JOB_TITLE  = "job_title"
	FIELD_LAST_NAME  = "last_name"
	FIELD_QUOTE      = "quote"
	FIELD_STATUS     = "status"
)

type Testimonial struct {
	dataobject.DataObject
}

func NewTestimonial() *Testimonial {
	return &Testimonial{}
}

func NewTestimonialFromEntity(store entitystore.StoreInterface, entity entitystore.EntityInterface) (*Testimonial, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}
	if entity.GetType() != ENTITY_TYPE {
		return nil, errors.New("invalid entity type")
	}

	// Fetch attributes from store
	attributes, err := store.EntityAttributeList(context.Background(), entity.GetID())
	if err != nil {
		return nil, err
	}

	testimonial := NewTestimonial()

	// Copy core entity fields
	testimonial.SetID(entity.GetID())
	testimonial.Set(entitystore.COLUMN_CREATED_AT, entity.GetCreatedAt())
	testimonial.Set(entitystore.COLUMN_UPDATED_AT, entity.GetUpdatedAt())

	for _, attr := range attributes {
		testimonial.Set(attr.GetKey(), attr.GetValue())
	}

	return testimonial, nil
}

// == SETTERS AND GETTERS =====================================================

func (t *Testimonial) Date() string {
	return t.Get(FIELD_DATE)
}

func (t *Testimonial) SetDate(date string) {
	t.Set(FIELD_DATE, date)
}

func (t *Testimonial) FirstName() string {
	return t.Get(FIELD_FIRST_NAME)
}

func (t *Testimonial) SetFirstName(firstName string) {
	t.Set(FIELD_FIRST_NAME, firstName)
}

func (t *Testimonial) ID() string {
	return t.Get(FIELD_ID)
}

func (t *Testimonial) SetID(id string) {
	t.Set(FIELD_ID, id)
}

func (t *Testimonial) ImageUrl() string {
	return t.Get(FIELD_IMAGE_URL)
}

func (t *Testimonial) SetImageUrl(imageUrl string) {
	t.Set(FIELD_IMAGE_URL, imageUrl)
}

func (t *Testimonial) JobTitle() string {
	return t.Get(FIELD_JOB_TITLE)
}

func (t *Testimonial) SetJobTitle(jobTitle string) {
	t.Set(FIELD_JOB_TITLE, jobTitle)
}

func (t *Testimonial) LastName() string {
	return t.Get(FIELD_LAST_NAME)
}

func (t *Testimonial) SetLastName(lastName string) {
	t.Set(FIELD_LAST_NAME, lastName)
}

func (t *Testimonial) Quote() string {
	return t.Get(FIELD_QUOTE)
}

func (t *Testimonial) SetQuote(quote string) {
	t.Set(FIELD_QUOTE, quote)
}

func (t *Testimonial) Status() string {
	return t.Get(FIELD_STATUS)
}

func (t *Testimonial) SetStatus(status string) {
	t.Set(FIELD_STATUS, status)
}

func (t *Testimonial) CreatedAt() string {
	return t.Get(entitystore.COLUMN_CREATED_AT)
}

func (t *Testimonial) SetCreatedAt(createdAt string) {
	t.Set(entitystore.COLUMN_CREATED_AT, createdAt)
}

func (t *Testimonial) UpdatedAt() string {
	return t.Get(entitystore.COLUMN_UPDATED_AT)
}

func (t *Testimonial) SetUpdatedAt(updatedAt string) {
	t.Set(entitystore.COLUMN_UPDATED_AT, updatedAt)
}
