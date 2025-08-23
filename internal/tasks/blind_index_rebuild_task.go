package tasks

import (
	"context"
	"errors"
	"project/internal/helpers"
	"project/internal/types"
	"slices"
	"strings"

	"github.com/dracory/blindindexstore"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/userstore"
)

// ============================================================================
// blindIndexRebuildTask
// ============================================================================
// Updates a blind index table, by upserting (insert, update) the data.
// If truncate is set to yes, the table will be truncated first, and then repopulated.
// If a data is empty it is useless in a search, it will not be inserted or updated,
// and will be removed if it exists.
// ============================================================================
// Example:
// - go run . task BlindIndexUpdate --index=all
// - go run . task BlindIndexUpdate --index=first_name
// - go run . task BlindIndexUpdate --index=first_name --truncate=yes
// ============================================================================
type blindIndexRebuildTask struct {
	taskstore.TaskHandlerBase

	app types.AppInterface

	// allowed indexes
	allowedIndexes []string

	// index to rebuild
	index string

	// truncate the index table
	truncate bool
}

var _ taskstore.TaskHandlerInterface = (*blindIndexRebuildTask)(nil) // verify it extends the task interface

// == CONSTRUCTOR =============================================================

func NewBlindIndexRebuildTask(app types.AppInterface) *blindIndexRebuildTask {
	return &blindIndexRebuildTask{
		app:            app,
		allowedIndexes: []string{"all", "email", "first_name", "last_name"},
	}
}

// == IMPLEMENTATION ==========================================================

func (task *blindIndexRebuildTask) Alias() string {
	return "BlindIndexUpdate"
}

func (task *blindIndexRebuildTask) Title() string {
	return "Blind Index Update"
}

func (task *blindIndexRebuildTask) Description() string {
	return "Truncates a blind index table, and repopulates it with the current data"
}

func (task *blindIndexRebuildTask) Enqueue(index string) (queuedTask taskstore.QueueInterface, err error) {
	if task.app.GetTaskStore() == nil {
		return nil, errors.New("task store is nil")
	}

	return task.app.GetTaskStore().TaskEnqueueByAlias(task.Alias(), map[string]any{
		"index": index,
	})
}

func (task *blindIndexRebuildTask) Handle() bool {
	task.index = task.GetParam("index")
	task.truncate = task.GetParam("truncate") == "yes"

	if !slices.Contains(task.allowedIndexes, task.index) {
		task.LogError("Invalid index: '" + task.index + "'. Must be one of: '" + strings.Join(task.allowedIndexes, "', '") + "'. Aborted.")
		return false
	}

	rebuilEmailIndex := task.index == "all" || task.index == "email"
	rebuildFirstNameIndex := task.index == "all" || task.index == "first_name"
	rebuildLastNameIndex := task.index == "all" || task.index == "last_name"

	if task.checkAndEnqueueTask() {
		return true
	}

	if rebuildFirstNameIndex && !task.rebuildFirstNameIndex(context.Background()) {
		return false
	}

	if rebuildLastNameIndex && !task.rebuildLastNameIndex(context.Background()) {
		return false
	}

	// !!! IMPORTANT
	// Set as last index to rebuild as its more important
	// than the other indexes
	// If the other indexes fail, nothing major will occur
	// But if this one fails, users will not be able to login
	// This is why let the other indexes rebuild first,
	// to make sure no errors have been spotted so far
	if rebuilEmailIndex && !task.rebuildEmailIndex(context.Background()) {
		return false
	}

	return true
}

func (task *blindIndexRebuildTask) rebuildEmailIndex(ctx context.Context) bool {
	task.LogInfo("Rebuilding email index:")

	if task.app.GetUserStore() == nil {
		task.LogError("User store is nil. Aborted.")
		return false
	}

	if task.truncate {
		task.LogInfo(" - Truncating blind index table...")
		if task.app.GetBlindIndexStoreEmail() == nil {
			task.LogError("BlindIndexStoreEmail is nil. Aborted.")
			return false
		}
		err := task.app.GetBlindIndexStoreEmail().Truncate()

		if err != nil {
			task.LogError("Error truncating blind index table: " + err.Error())
			return false
		}
	}

	task.LogInfo(" - Fetching users list...")
	users, err := task.app.GetUserStore().UserList(ctx, userstore.NewUserQuery())

	if err != nil {
		task.LogError("Error retrieving users: " + err.Error())
		return false
	}

	task.LogInfo(" - Rebuilding index...")
	for _, user := range users {
		if !task.insertEmailForUser(ctx, user) {
			task.LogError("- Failed to insert email for user: " + user.ID() + ". Aborted.")
			return false
		}
	}

	task.LogInfo(" - Index rebuilt successfully")
	return true
}

func (task *blindIndexRebuildTask) rebuildFirstNameIndex(ctx context.Context) bool {
	task.LogInfo("Rebuilding first name index:")

	if task.app.GetUserStore() == nil {
		task.LogError("User store is nil. Aborted.")
		return false
	}

	if task.truncate {
		task.LogInfo(" - Truncating blind index table")
		if task.app.GetBlindIndexStoreFirstName() == nil {
			task.LogError("BlindIndexStoreFirstName is nil. Aborted.")
			return false
		}
		err := task.app.GetBlindIndexStoreFirstName().Truncate()

		if err != nil {
			task.LogError("Error truncating blind index table: " + err.Error())
			return false
		}
	}

	task.LogInfo(" - Fetching users list")
	users, err := task.app.GetUserStore().UserList(ctx, userstore.NewUserQuery())

	if err != nil {
		task.LogError("Error retrieving users: " + err.Error())
		return false
	}

	task.LogInfo(" - Rebuilding index")
	for _, user := range users {
		if !task.insertFirstNameForUser(ctx, user) {
			task.LogError("- Failed to insert first name for user: " + user.ID() + ". Aborted.")
			return false
		}
	}

	task.LogInfo(" - Index rebuilt successfully")
	return true
}

func (task *blindIndexRebuildTask) rebuildLastNameIndex(ctx context.Context) bool {
	task.LogInfo("Rebuilding last name index:")

	if task.app.GetUserStore() == nil {
		task.LogError("User store is nil. Aborted.")
		return false
	}

	if task.truncate {
		task.LogInfo(" - Truncating blind index table")
		if task.app.GetBlindIndexStoreLastName() == nil {
			task.LogError("BlindIndexStoreLastName is nil. Aborted.")
			return false
		}
		err := task.app.GetBlindIndexStoreLastName().Truncate()

		if err != nil {
			task.LogError("Error truncating blind index table: " + err.Error())
			return false
		}
	}

	task.LogInfo(" - Fetching users list")
	users, err := task.app.GetUserStore().UserList(ctx, userstore.NewUserQuery())

	if err != nil {
		task.LogError("Error retrieving users: " + err.Error())
		return false
	}

	task.LogInfo(" - Rebuilding index")
	for _, user := range users {
		if !task.insertLastNameForUser(ctx, user) {
			task.LogError("- Failed to insert last name for user: " + user.ID() + ". Aborted.")
			return false
		}
	}

	task.LogInfo(" - Index rebuilt successfully")
	return true
}

func (task *blindIndexRebuildTask) insertEmailForUser(ctx context.Context, user userstore.UserInterface) bool {
	searchValue, err := task.app.GetBlindIndexStoreEmail().SearchValueFindBySourceReferenceID(user.ID())

	if err != nil {
		task.LogError("Error searching for blind index by source reference ID: " + user.ID() + " - " + err.Error())
		return false
	}

	isIndexed := searchValue != nil

	emailToken := user.Email()

	// No need to index an empty email
	if emailToken == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreEmail().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty email, nothing to do
	}

	m, err := helpers.Untokenize(ctx, task.app.GetVaultStore(), task.app.GetConfig().GetVaultKey(), map[string]string{"email": emailToken})

	if err != nil {
		task.LogError("Error untokenizing user token: " + emailToken + " - " + err.Error())
		return false
	}

	email := m["email"]

	if email == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreEmail().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty email, nothing to do
	}

	// Upsert
	if isIndexed {
		searchValue.SetSearchValue(email)
		err = task.app.GetBlindIndexStoreEmail().SearchValueUpdate(searchValue)
		if err != nil {
			task.LogError("Error updating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	} else {
		err = task.app.GetBlindIndexStoreEmail().SearchValueCreate(blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.ID()).
			SetSearchValue(email))
		if err != nil {
			task.LogError("Error creating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	}

	return true
}

func (task *blindIndexRebuildTask) insertFirstNameForUser(ctx context.Context, user userstore.UserInterface) bool {
	if task.app.GetBlindIndexStoreFirstName() == nil {
		task.LogError("BlindIndexStoreFirstName is nil. Aborted.")
		return false
	}
	searchValue, err := task.app.GetBlindIndexStoreFirstName().SearchValueFindBySourceReferenceID(user.ID())

	if err != nil {
		task.LogError("Error searching for blind index by source reference ID: " + user.ID() + " - " + err.Error())
		return false
	}

	isIndexed := searchValue != nil

	firstNameToken := user.FirstName()

	// No need to index an empty first name
	if firstNameToken == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreFirstName().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty first name, nothing to do
	}

	m, err := helpers.Untokenize(ctx, task.app.GetVaultStore(), task.app.GetConfig().GetVaultKey(), map[string]string{"first_name": firstNameToken})

	if err != nil {
		task.LogError("Error untokenizing user token: " + firstNameToken + " - " + err.Error())
		return false
	}

	firstName := m["first_name"]

	// No need to index an empty first name
	if firstName == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreFirstName().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty first name, nothing to do
	}

	// Upsert the search value
	if isIndexed {
		searchValue.SetSearchValue(firstName)
		err = task.app.GetBlindIndexStoreFirstName().SearchValueUpdate(searchValue)
		if err != nil {
			task.LogError("Error updating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	} else {
		err = task.app.GetBlindIndexStoreFirstName().SearchValueCreate(blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.ID()).
			SetSearchValue(firstName))
		if err != nil {
			task.LogError("Error creating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	}

	return true
}

func (task *blindIndexRebuildTask) insertLastNameForUser(ctx context.Context, user userstore.UserInterface) bool {
	if task.app.GetBlindIndexStoreLastName() == nil {
		task.LogError("BlindIndexStoreLastName is nil. Aborted.")
		return false
	}
	searchValue, err := task.app.GetBlindIndexStoreLastName().SearchValueFindBySourceReferenceID(user.ID())

	if err != nil {
		task.LogError("Error searching for blind index by source reference ID: " + user.ID() + " - " + err.Error())
		return false
	}

	isIndexed := searchValue != nil

	lastNameToken := user.LastName()

	// No need to index an empty last name
	if lastNameToken == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreLastName().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty last name, nothing to do
	}

	m, err := helpers.Untokenize(ctx, task.app.GetVaultStore(), task.app.GetConfig().GetVaultKey(), map[string]string{"last_name": lastNameToken})

	if err != nil {
		task.LogError("Error untokenizing user token: " + lastNameToken + " - " + err.Error())
		return false
	}

	lastName := m["last_name"]

	// No need to index an empty last name
	if lastName == "" {
		if isIndexed {
			err = task.app.GetBlindIndexStoreLastName().SearchValueDelete(searchValue)
			if err != nil {
				task.LogError("Error deleting blind index for user: " + user.ID() + " - " + err.Error())
			}
		}
		return true // empty last name, nothing to do
	}

	// Upsert the search value
	if isIndexed {
		searchValue.SetSearchValue(lastName)
		err = task.app.GetBlindIndexStoreLastName().SearchValueUpdate(searchValue)
		if err != nil {
			task.LogError("Error updating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	} else {
		err = task.app.GetBlindIndexStoreLastName().SearchValueCreate(blindindexstore.NewSearchValue().
			SetSourceReferenceID(user.ID()).
			SetSearchValue(lastName))
		if err != nil {
			task.LogError("Error creating blind index for user: " + user.ID() + " - " + err.Error())
			return false
		}
	}

	return true
}

func (task *blindIndexRebuildTask) checkAndEnqueueTask() bool {
	// ... (rest of the code remains the same)
	if task.HasQueuedTask() {
		return false
	}

	// 2. Is the task asked to be enqueued?
	if task.GetParam("enqueue") != "yes" {
		return false
	}

	// 3. Enqueue the task
	_, err := task.Enqueue(task.index)

	if err != nil {
		task.LogError("Error enqueuing task: " + err.Error())
	} else {
		task.LogSuccess("Task enqueued.")
	}

	return true
}
