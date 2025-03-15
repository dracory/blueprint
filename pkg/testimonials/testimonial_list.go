package testimonials

import (
	"github.com/gouniverse/entitystore"
)

func TestimonialList(entityStore *entitystore.Store) ([]Testimonial, error) {
	result, err := entityStore.EntityList(entitystore.EntityQueryOptions{
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
