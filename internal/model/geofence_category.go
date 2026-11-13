package model

type GeofenceCategory string

const (
	DeliveryZone   GeofenceCategory = "delivery_zone"
	RestrictedZone GeofenceCategory = "restricted_zone"
	TollZone       GeofenceCategory = "toll_zone"
	CustomerArea   GeofenceCategory = "customer_area"
)

func (c GeofenceCategory) IsValid() bool {
	switch c {
	case DeliveryZone,
		RestrictedZone,
		TollZone,
		CustomerArea:
		return true
	default:
		return false
	}
}