package parse

import "fmt"

type Condition struct {
	Key   string
	Op    string
	Value string
}

type AndGroup []Condition

type Groups []AndGroup

type Conditions map[string]Groups

func Query(qs []interface{}) Conditions {
	c := make(Conditions)

	for _, q := range qs {
		mq, ok := q.(map[string]interface{})
		if !ok {
			continue
		}

		for k, v := range mq {
			nv, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			doQuery(c, k, nv)
		}
	}

	// fmt.Printf("Result: %#v\n", c)
	return c
}

func doQuery(c Conditions, pk string, q map[string]interface{}) {
	for k, v := range q {
		nv, ok := v.(map[string]interface{})

		if ok {
			doQuery(c, k, nv)

		} else {
			g, ok := c[pk]
			if !ok {
				g = Groups{}
			}

			ag := AndGroup{}

			for k, v := range q {
				ag = append(ag, Condition{
					Key:   pk,
					Op:    k,
					Value: fetchString(v),
				})
			}

			g = append(g, ag)

			c[pk] = g

			break
		}
	}
}

func fetchString(v interface{}) string {
	switch v := v.(type) {
	case *string:
		return *v
	case *int8:
		return fmt.Sprintf("%d", *v)
	case *int:
		return fmt.Sprintf("%d", *v)
	case *int32:
		return fmt.Sprintf("%d", *v)
	case *int64:
		return fmt.Sprintf("%d", *v)
	case *float32:
		return fmt.Sprintf("%f", *v)
	case *float64:
		return fmt.Sprintf("%f", *v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%f", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case []byte:
		return string(v)
	case nil:
		return "0"
	default:
		return fmt.Sprintf("%s", v)
	}
}
