package crud 

type Service[E Entity] struct {
  Store Store[E] 
  IdPolicy IdPolicy
  ReadP, WriteP PermissionPolicy 
}

type PermissionPolicy = func(RequestData) error
type IdPolicy = func(RequestData) (id string, err error)

type RequestData struct {
  IdFromURL string 
  CallerId string 
}
 
func (s Service[E]) Create(rd RequestData, e E) error {
  if err := s.WriteP(rd); err != nil {
    return err
  }
  return s.Store.Create(e)
} 
func (s Service[E]) Read(rd RequestData) (e E, err error) {
  if err := s.ReadP(rd); err != nil {
    return e, err
  }
  id, err := s.IdPolicy(rd)
  if err != nil {
    return e, err
  }
  return s.Store.Read(id)
}
func (s Service[E]) Update(rd RequestData, e E) error {
  if err := s.WriteP(rd); err != nil {
    return err
  }
  id, err := s.IdPolicy(rd) 
  if err != nil {
    return err
  }
  return s.Store.Update(id, e)
} 
func (s Service[E]) Delete(rd RequestData) error {
  if err := s.WriteP(rd); err != nil {
    return err
  }
  id, err := s.IdPolicy(rd) 
  if err != nil {
    return err
  }
  return s.Store.Delete(id)
} 



