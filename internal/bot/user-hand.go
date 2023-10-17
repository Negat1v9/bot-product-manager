package telegram

// tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Info: create message with information to user, about he have to forwad its
// func (h *Hub) answerToCreateList(ChatID int64) *tg.MessageConfig {
// 	msg := h.createMessage(ChatID, answerCreateListMsg)
// 	return msg
// }

// // Info: Create list in database, and answer with new list name
// func (h *Hub) createList(ChatID int64, list *store.ProductList) (*tg.MessageConfig, error) {

// 	err := h.db.ProductList().Create(context.TODO(), list)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", *list.Name))
// 	return msg, nil
// }

// // Info: Select all users lists and create inline keyboard with its
// func (h *Hub) getListName(UserID, ChatID int64) (msg *tg.MessageConfig, err error) {
// 	lists, err := h.db.ProductList().GetAll(context.TODO(), int(UserID))
// 	if err != nil {
// 		if err == store.NoRowListOfProductError {
// 			msg = h.createMessage(ChatID, "Nothing is found. Create Youre First list!")
// 			return msg, nil
// 		}
// 		return nil, err
// 	}
// 	keyboard := createListProductInline(lists)
// 	msg = h.createMessage(ChatID, "List of Product-lists")
// 	msg.ReplyMarkup = keyboard
// 	return msg, nil
// }

// func (h *Hub) getProductList(ChatID int64, listID int, listName string) (msg *tg.MessageConfig, err error) {
// 	product, err := h.db.Product().GetAll(context.TODO(), listID)
// 	if err != nil {
// 		if err == store.NoRowProductError {

// 			msg = h.createMessage(ChatID, emptyListMessage)
// 			msg.ReplyMarkup = createProductsInline(listName)
// 			return msg, nil

// 		}
// 		return nil, err
// 	}
// 	text := createMessageProductList(product.Products)
// 	msg = h.createMessage(ChatID, text)

// 	msg.ReplyMarkup = createProductsInline(listName)
// 	return msg, nil
// }

// func (h *Hub) addNewProduct(ChatID int64, Products, listName string) (*tg.MessageConfig, error) {
// 	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	newProduct := parseStringToProducts(Products, listID)

// 	lastProduct, err := h.db.Product().GetAll(context.TODO(), listID)
// 	if err != nil {
// 		// if not row exist
// 		if err == store.NoRowProductError {
// 			err = h.db.Product().Create(context.TODO(), listID)
// 			lastProduct = &store.Product{
// 				ListID:   listID,
// 				Products: []string{},
// 			}
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			return nil, err
// 		}
// 	}
// 	text := createMessageSuccessAddedProduct(newProduct.Products)
// 	// add to old values new
// 	if len(lastProduct.Products) > 0 {
// 		lastProduct.Products = append(lastProduct.Products, newProduct.Products...)
// 		newProduct.Products = lastProduct.Products
// 	}
// 	if err := h.db.Product().Add(context.TODO(), newProduct); err != nil {
// 		return nil, err
// 	}

// 	msg := h.createMessage(ChatID, text)
// 	msg.ReplyMarkup = createInlineGetCurList(listID, listName)
// 	return msg, nil
// }

// func (h *Hub) createMessageForEditList(ChatID int64, listName string) *tg.MessageConfig {
// 	msg := h.createMessage(ChatID, answerEditListMessage+listName)
// 	return msg
// }

// func (h *Hub) compliteProductList(ChatID int64, productListName string) (*tg.MessageConfig, error) {
// 	listID, err := h.db.ProductList().GetListID(context.TODO(), productListName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = h.db.ProductList().Delete(context.TODO(), listID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(ChatID, isCompletesProductListMsg+productListName)
// 	return msg, nil
// }

// func (h *Hub) editProductList(chatID int64, listName string, indexProducts map[int]bool) (*tg.MessageConfig, error) {
// 	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	products, err := h.db.Product().GetAll(context.TODO(), listID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	products.Products = deleteProductByIndex(products.Products, indexProducts)

// 	err = h.db.Product().Add(context.TODO(), *products)
// 	if err != nil {
// 		// TODO: Not deleted
// 		return nil, err
// 	}
// 	text := createMessageProductList(products.Products)
// 	msg := h.createMessage(chatID, text)
// 	msg.ReplyMarkup = createProductsInline(listName)
// 	return msg, nil
// }

// func (h *Hub) GetGroupLists(UserID int64, groupID int) (msg *tg.MessageConfig, err error) {
// 	var isOwnerGroup = false
// 	groupList, err := h.db.ManagerGroup().AllByGroupID(context.TODO(), groupID)
// 	if err != nil {
// 		if err == store.NoRowListOfProductError {
// 			msg = h.createMessage(UserID, err.Error())
// 		} else {
// 			return nil, err
// 		}
// 	} else {
// 		msg = h.createMessage(UserID, "group lists")
// 	}
// 	if UserID == groupList.GroupOwnerID {
// 		isOwnerGroup = true
// 	}
// 	msg.ReplyMarkup = createInlineGroupList(groupList.PruductLists, groupID, isOwnerGroup)
// 	return msg, nil
// }

// func (h *Hub) createMessageCreateGroupList(UserID int64, groupID int) (*tg.MessageConfig, error) {
// 	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(UserID, answerCreateGroupListMsg+group.GroupName)
// 	return msg, nil
// }

// func (h *Hub) createGroupList(UserID int64, listName, groupName string) (*tg.MessageConfig, error) {
// 	group, err := h.db.ManagerGroup().ByGroupName(context.TODO(), groupName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	list := &store.ProductList{
// 		OwnerID: &UserID,
// 		GroupID: &group.ID,
// 		Name:    &listName,
// 	}
// 	err = h.db.ProductList().Create(context.TODO(), list)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(UserID, `New list is created`)
// 	return msg, nil
// }

// func (h *Hub) getUserForDeleteFrGr(ChatID int64, groupID int) (msg *tg.MessageConfig, err error) {
// 	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Info: if user at group lowwer them 2, means at group only owner
// 	if len(*groupInfo.UsersInfo) < 2 {
// 		msg = h.createMessage(ChatID, emptyUserInGroup)
// 		return
// 	}
// 	msg = h.createMessage(ChatID, "Choise user from:"+groupInfo.GroupName)
// 	msg.ReplyMarkup = createInlineDeleteUser(*groupInfo.UsersInfo, groupID)
// 	return
// }

// func (h *Hub) deleteUserFromGroup(ChatID, userID int64, groupID int) (*tg.MessageConfig, error) {
// 	g := &store.Group{
// 		UserID:  userID,
// 		GroupID: groupID,
// 	}
// 	err := h.db.Group().DeleteUser(context.TODO(), g)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(ChatID, successDeletedUser)
// 	return msg, nil
// }

// func (h *Hub) createMessageForInviteUser(chatID int64, groupID int) (*tg.MessageConfig, error) {
// 	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(chatID, textForInvitingNewUser+group.GroupName)
// 	return msg, nil
// }

// func (h *Hub) inviteNewUser(ChatID int64, newUserName, groupName string) (*tg.MessageConfig, error) {
// 	group, err := h.db.ManagerGroup().ByGroupName(context.TODO(), groupName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = h.sendInviteMessage(newUserName, group.ID, ChatID)
// 	if err != nil {
// 		if err == userAlredyGroupError || err == userNoExistError {
// 			msg := h.createMessage(ChatID, err.Error())
// 			return msg, nil
// 		}
// 	}
// 	msg := h.createMessage(ChatID, inviteSendMessage)
// 	return msg, nil
// }

// func (h *Hub) userReadyJoinGroup(ChatID, newUserID int64, groupID int) (*tg.MessageConfig, error) {
// 	g := &store.Group{
// 		UserID:  newUserID,
// 		GroupID: groupID,
// 	}
// 	err := h.db.Group().AddUser(context.TODO(), g)
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg := h.createMessage(ChatID, userInvitedInGroupMessage)
// 	return msg, nil
// }

// // func (h *Hub) userRefuseJoinGroup()
