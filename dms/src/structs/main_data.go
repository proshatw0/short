package structs

type Link_source struct {
	Original_link Hash_Table
	Short_link    Hash_Table
}


func NewLinkSource(len int) Link_source {
	original_link := NewHashTable(len)
	short_link := NewHashTable(len)
 
	return Link_source{
	   Original_link: original_link,
	   Short_link:    short_link,
	}
}

