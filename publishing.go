package wecms

import "errors"

type Publishing struct {
	srcRep    *Repository
	targetRep *Repository
	currUser  *User
}

func (pub *Publishing) PublishItem(item *Item) error {
	if pub.srcRep == nil {
		return errors.New("Invalid source repository: the srcRep cannot be nil")
	}
	if pub.targetRep == nil {
		return errors.New("Invalid target repository: the targetRep cannot be nil")
	}
	if item == nil {
		return errParamNil("item")
	}
	if item.currentRep != pub.srcRep {
		return errors.New("This item cannot be published. The item repository is not the same as the srcRep")
	}
	if pub.srcRep == pub.targetRep {
		return nil
	}
	tmpItem := pub.targetRep.itemCache[item.Id]
	delete(pub.targetRep.itemCache, item.Id)
	editing := pub.targetRep.Editing(pub.currUser)
	pubItem := *item
	pubItem.currentRep = pub.targetRep
	err := editing.SaveItem(&pubItem)
	if err != nil {
		pub.targetRep.itemCache[item.Id] = tmpItem
		return err
	}
	return nil
}