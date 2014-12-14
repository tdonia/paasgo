package paas

import (
	"appengine/datastore"
)

type Ancestor struct {
	Context   Context
	Kind      string
	KeyString string
	KeyInt    int64
	Parent    *datastore.Key
}

type Query struct {
	Kind      string
	Context   Context
	Offset    int
	Limit     int
	KeyString string
	KeyInt    int64
	Ancestors []Ancestor
	Order     string
	Filters   map[string]string
}

func (query Query) DeleteByKey(key *datastore.Key) error {
	return datastore.Delete(query.Context, key)
}

func (query Query) Delete() error {
	key := query.Key()
	return query.DeleteByKey(key)
}

func (query Query) Put(entry interface{}) error {
	key := query.Key()
	_, err := datastore.Put(query.Context, key, entry)

	if err != nil {
		query.Context.Errorf("Put Query Error Kind: " + query.Kind)
		query.Context.Errorf("Put Query Error KeyString: " + query.KeyString)
		for _, anc := range query.Ancestors {
			query.Context.Errorf("Put Query Error-anc Kind: " + anc.Kind)
			query.Context.Errorf("Put Query Error-anc KeyString: " + anc.KeyString)
		}
		query.Context.Errorf(err.Error())
	}
	return err
}

func (query Query) Get(result interface{}) error {
	key := query.Key()
	return datastore.Get(query.Context, key, result)
}

func (query Query) GetAll(results interface{}) error {
	_, err := query.CreateQuery().GetAll(query.Context, results)
	return err
}

func (query Query) AncestorKey() (parent *datastore.Key) {
	if len(query.Ancestors) > 0 {
		for _, ancestor := range query.Ancestors {
			ancestor.Context = query.Context
			ancestor.Parent = parent
			parent = ancestor.Key()
		}
		return parent
	}
	return nil
}

func (ancestor Ancestor) Key() *datastore.Key {
	return datastore.NewKey(ancestor.Context, ancestor.Kind, ancestor.KeyString, ancestor.KeyInt, ancestor.Parent)
}

func (query Query) Key() *datastore.Key {
	ancestor_key := query.AncestorKey()
	return datastore.NewKey(query.Context, query.Kind, query.KeyString, query.KeyInt, ancestor_key)
}

func (query Query) CreateQuery() (q *datastore.Query) {

	q = datastore.NewQuery(query.Kind)

	if len(query.Ancestors) > 0 {
		ancestor_key := query.AncestorKey()
		q = q.Ancestor(ancestor_key)
	}

	for filter_by, value := range query.Filters {
		q = q.Filter(filter_by, value)
	}

	if query.Limit != 0 {
		q = q.Limit(query.Limit)
	}

	if query.Offset != 0 {
		q = q.Offset(query.Offset)
	}

	if query.Order != "" {
		q = q.Order(query.Order)
	}

	return q
}
