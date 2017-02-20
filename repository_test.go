package wecms

import "testing"

func initRep(t *testing.T) error {
	rep,err := NewRepository("master", "mongodb://localhost", "master", 100,1000)
	if err != nil {
		println(err.Error())
		t.Failed()
		return err
	}
	err = RegRepository(rep)
	if err != nil {
		println(err.Error())
		t.Failed()
		return err
	}
	return nil
}

var testTemplateId ID = "fab42f1e-3cec-4855-a8c7-b283d8e12498"

var testAdmin = &User{
	Id: "50408a92-bd35-48c6-814f-b970f6cedfaa",
	Email: "test@test.com",
	UserName: "Testing",
	FullName: "Testing User",
	FirstName: "Testing",
	LastName: "User",
	Roles: []RoleType {
		Administrator,
	},
}

func TestRepository_SaveTemplate(t *testing.T) {
	if err := initRep(t); err != nil {
		t.Fatal(err.Error())
		return
	}
	newRep := GetRepository("master").Editing(testAdmin)
	if newRep == nil {
		t.Fatal("cannot find the editing by name 'master'")
		return
	}
	temp := &Template{
		Id: testTemplateId,
		Name: "TestTemplate",
		Sections: []*TemplateSection{
			{
				Name: "Content",
				Fields: []*TemplateField{
					{
						Id:              NewID(),
						Name:            "Name",
						DisplayTitle:    "Name",
						FieldType:       "Single-Line Text",
						Mandatory:       true,
						ValidationRegex: "",
						DefaultValue:    "",
					},
					{
						Id:              NewID(),
						Name:            "Body Content",
						DisplayTitle:    "The Body Content",
						FieldType:       "Rich Text",
						Mandatory:       false,
						ValidationRegex: "",
						DefaultValue:    "",
					},
				},
			},
		},
	}
	err := newRep.saveTemplate(temp)
	if err != nil {
		println(err.Error())
		t.Failed()
	}
}

func TestRepository_GetTemplate(t *testing.T) {
	if err := initRep(t);err != nil {
		return
	}
	rep := GetRepository("master")
	tmp := rep.GetTemplate(testTemplateId)
	if tmp == nil {
		t.Fatal("cannot find template by id:" + string(testTemplateId))
	} else {
		println(tmp.GetField("Body Content").Id)
	}
}