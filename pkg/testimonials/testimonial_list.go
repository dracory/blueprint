package testimonials

import (
	"github.com/dracory/entitystore"
)

func TestimonialList(store entitystore.StoreInterface) ([]Testimonial, error) {
	result, err := store.EntityList(entitystore.EntityQueryOptions{
		EntityType: "testimonial",
	})

	if err != nil {
		return nil, err
	}

	testimonials := []Testimonial{}

	for _, entry := range result {
		testimonial, err := NewTestimonialFromEntity(entry)

		if err != nil {
			continue
		}

		testimonials = append(testimonials, *testimonial)
	}

	return testimonials, nil
}
