package tests

type FakeEvent struct {
}

func (f FakeEvent) ID() string {
    //TODO implement me
    panic("implement me")
}

func (f FakeEvent) Type() string {
    //TODO implement me
    panic("implement me")
}

func (f FakeEvent) DataContentType() string {
    //TODO implement me
    panic("implement me")
}

func (f FakeEvent) Serialize() ([]byte, error) {
    //TODO implement me
    panic("implement me")
}

func (f FakeEvent) Accept(visitor FakeEventVisitor) error {
    return visitor.VisitFakeEvent(f)
}

type FakeEventVisitor interface {
    VisitFakeEvent(e FakeEvent) error
}
