package conversation

import (
	"fmt"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/gou/connector"
	"github.com/yaoapp/xun"
	"github.com/yaoapp/xun/capsule"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/test"
)

func TestNewXunDefault(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")

	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}

	err = capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}

	err = capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})

	if err != nil {
		t.Error(err)
		return
	}

	// Check history table
	has, err := capsule.Schema().HasTable("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// Check chat table
	has, err = capsule.Schema().HasTable("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// Check assistant table
	has, err = capsule.Schema().HasTable("__unit_test_conversation_assistant")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// validate the history table
	tab, err := conv.schema.GetTable(conv.getHistoryTable())
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"id", "sid", "cid", "uid", "role", "name", "content", "context", "created_at", "updated_at", "expired_at"}
	for _, field := range fields {
		assert.Equal(t, true, tab.HasColumn(field))
	}

	// validate the chat table
	tab, err = conv.schema.GetTable(conv.getChatTable())
	if err != nil {
		t.Fatal(err)
	}

	chatFields := []string{"id", "chat_id", "title", "sid", "created_at", "updated_at"}
	for _, field := range chatFields {
		assert.Equal(t, true, tab.HasColumn(field))
	}

	// validate the assistant table
	tab, err = conv.schema.GetTable(conv.getAssistantTable())
	if err != nil {
		t.Fatal(err)
	}

	assistantFields := []string{"id", "assistant_id", "type", "name", "avatar", "connector", "description", "options", "prompts", "flows", "files", "functions", "tags", "readonly", "permissions", "automated", "mentionable", "created_at", "updated_at"}
	for _, field := range assistantFields {
		assert.Equal(t, true, tab.HasColumn(field))
	}
}

func TestNewXunConnector(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()

	conn, err := connector.Select("mysql")
	if err != nil {
		t.Fatal(err)
	}

	sch, err := conn.Schema()
	if err != nil {
		t.Fatal(err)
	}

	defer sch.DropTableIfExists("__unit_test_conversation_history")
	defer sch.DropTableIfExists("__unit_test_conversation_chat")
	defer sch.DropTableIfExists("__unit_test_conversation_assistant")

	sch.DropTableIfExists("__unit_test_conversation_history")
	sch.DropTableIfExists("__unit_test_conversation_chat")
	sch.DropTableIfExists("__unit_test_conversation_assistant")

	conv, err := NewXun(Setting{
		Connector: "mysql",
		Table:     "__unit_test_conversation",
	})

	if err != nil {
		t.Error(err)
		return
	}

	// Check history table
	has, err := sch.HasTable("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// Check chat table
	has, err = sch.HasTable("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// Check assistant table
	has, err = sch.HasTable("__unit_test_conversation_assistant")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, has)

	// validate the history table
	tab, err := conv.schema.GetTable(conv.getHistoryTable())
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"id", "sid", "cid", "uid", "role", "name", "content", "context", "created_at", "updated_at", "expired_at"}
	for _, field := range fields {
		assert.Equal(t, true, tab.HasColumn(field))
	}

	// validate the chat table
	tab, err = conv.schema.GetTable(conv.getChatTable())
	if err != nil {
		t.Fatal(err)
	}

	chatFields := []string{"id", "chat_id", "title", "sid", "created_at", "updated_at"}
	for _, field := range chatFields {
		assert.Equal(t, true, tab.HasColumn(field))
	}

	// validate the assistant table
	tab, err = conv.schema.GetTable(conv.getAssistantTable())
	if err != nil {
		t.Fatal(err)
	}

	assistantFields := []string{"id", "assistant_id", "type", "name", "avatar", "connector", "description", "options", "prompts", "flows", "files", "functions", "tags", "readonly", "permissions", "automated", "mentionable", "created_at", "updated_at"}
	for _, field := range assistantFields {
		assert.Equal(t, true, tab.HasColumn(field))
	}
}

func TestXunSaveAndGetHistory(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")

	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}

	err = capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
		TTL:       3600,
	})

	// save the history
	cid := "123456"
	err = conv.SaveHistory("123456", []map[string]interface{}{
		{"role": "user", "name": "user1", "content": "hello"},
		{"role": "assistant", "name": "user1", "content": "Hello there, how"},
	}, cid, nil)
	assert.Nil(t, err)

	// get the history
	data, err := conv.GetHistory("123456", cid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(data))
}

func TestXunSaveAndGetHistoryWithCID(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")

	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}

	err = capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
		TTL:       3600,
	})

	// save the history with specific cid
	sid := "123456"
	cid := "789012"
	messages := []map[string]interface{}{
		{"role": "user", "name": "user1", "content": "hello"},
		{"role": "assistant", "name": "assistant1", "content": "Hi! How can I help you?"},
	}
	err = conv.SaveHistory(sid, messages, cid, nil)
	assert.Nil(t, err)

	// get the history for specific cid
	data, err := conv.GetHistory(sid, cid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(data))

	// save another message with different cid
	anotherCID := "345678"
	moreMessages := []map[string]interface{}{
		{"role": "user", "name": "user1", "content": "another message"},
	}
	err = conv.SaveHistory(sid, moreMessages, anotherCID, nil)
	assert.Nil(t, err)

	// get history for the first cid - should still be 2 messages
	data, err = conv.GetHistory(sid, cid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(data))

	// get history for the second cid - should be 1 message
	data, err = conv.GetHistory(sid, anotherCID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(data))

	// get all history for the sid without specifying cid
	allData, err := conv.GetHistory(sid, cid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(allData))
}

func TestXunGetChats(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")

	// Drop both tables before test
	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	if err != nil {
		t.Fatal(err)
	}
	err = capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Save some test chats
	sid := "test_user"
	messages := []map[string]interface{}{
		{"role": "user", "content": "test message"},
	}

	// Create chats with different dates
	for i := 0; i < 5; i++ {
		chatID := fmt.Sprintf("chat_%d", i)
		// First create the chat with a title
		err = conv.newQueryChat().Insert(map[string]interface{}{
			"chat_id":    chatID,
			"title":      fmt.Sprintf("Test Chat %d", i),
			"sid":        sid,
			"created_at": time.Now(),
		})
		if err != nil {
			t.Fatal(err)
		}

		// Then save the history
		err = conv.SaveHistory(sid, messages, chatID, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test getting chats with default filter
	filter := ChatFilter{
		PageSize: 10,
		Order:    "desc",
	}
	groups, err := conv.GetChats(sid, filter)
	if err != nil {
		t.Fatal(err)
	}

	assert.Greater(t, len(groups.Groups), 0)

	// Test with keywords
	filter.Keywords = "test"
	groups, err = conv.GetChats(sid, filter)
	if err != nil {
		t.Fatal(err)
	}

	assert.Greater(t, len(groups.Groups), 0)
}

func TestXunDeleteChat(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a test chat
	sid := "test_user"
	cid := "test_chat"
	messages := []map[string]interface{}{
		{"role": "user", "content": "test message"},
	}

	// Save the chat and history
	err = conv.SaveHistory(sid, messages, cid, nil)
	assert.Nil(t, err)

	// Verify chat exists
	chat, err := conv.GetChat(sid, cid)
	assert.Nil(t, err)
	assert.NotNil(t, chat)

	// Delete the chat
	err = conv.DeleteChat(sid, cid)
	assert.Nil(t, err)

	// Verify chat is deleted
	chat, err = conv.GetChat(sid, cid)
	assert.Nil(t, err)
	assert.Equal(t, (*ChatInfo)(nil), chat)
}

func TestXunDeleteAllChats(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_chat")

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create multiple test chats
	sid := "test_user"
	messages := []map[string]interface{}{
		{"role": "user", "content": "test message"},
	}

	// Save multiple chats
	for i := 0; i < 3; i++ {
		cid := fmt.Sprintf("test_chat_%d", i)
		err = conv.SaveHistory(sid, messages, cid, nil)
		assert.Nil(t, err)
	}

	// Verify chats exist
	response, err := conv.GetChats(sid, ChatFilter{})
	assert.Nil(t, err)
	assert.Greater(t, response.Total, int64(0))

	// Delete all chats
	err = conv.DeleteAllChats(sid)
	assert.Nil(t, err)

	// Verify all chats are deleted
	response, err = conv.GetChats(sid, ChatFilter{})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), response.Total)
}

func TestXunAssistantCRUD(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")

	// Drop assistant table before test
	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test creating a new assistant
	tagsJSON, err := jsoniter.MarshalToString([]string{"tag1", "tag2", "tag3"})
	if err != nil {
		t.Fatal(err)
	}

	optionsJSON, err := jsoniter.MarshalToString(map[string]interface{}{
		"model": "gpt-4",
	})
	if err != nil {
		t.Fatal(err)
	}

	assistant := map[string]interface{}{
		"name":        "Test Assistant",
		"type":        "assistant",
		"avatar":      "https://example.com/avatar.png",
		"connector":   "openai",
		"description": "Test Description",
		"tags":        tagsJSON,
		"options":     optionsJSON,
	}

	// Test SaveAssistant (Create)
	err = conv.SaveAssistant(assistant)
	assert.Nil(t, err)
	assistantID := assistant["assistant_id"].(string)
	assert.NotEmpty(t, assistantID)

	// Test GetAssistants with no filter
	resp, err := conv.GetAssistants(AssistantFilter{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.P.Items))

	// Test GetAssistants with tag filter (single tag)
	resp, err = conv.GetAssistants(AssistantFilter{
		Tags: []string{"tag1"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.P.Items))

	// Test GetAssistants with tag filter (multiple tags)
	resp, err = conv.GetAssistants(AssistantFilter{
		Tags: []string{"tag1", "tag4"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.P.Items))

	// Test GetAssistants with non-existent tag
	resp, err = conv.GetAssistants(AssistantFilter{
		Tags: []string{"nonexistent"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp.P.Items))

	// Test SaveAssistant (Update)
	assistant["name"] = "Updated Assistant"
	err = conv.SaveAssistant(assistant)
	assert.Nil(t, err)

	resp, err = conv.GetAssistants(AssistantFilter{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.P.Items))
	item := resp.P.Items[0].(xun.R)
	assert.Equal(t, "Updated Assistant", item["name"])

	// Test DeleteAssistant
	err = conv.DeleteAssistant(assistantID)
	assert.Nil(t, err)

	resp, err = conv.GetAssistants(AssistantFilter{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp.P.Items))
}

func TestXunAssistantPagination(t *testing.T) {
	test.Prepare(t, config.Conf)
	defer test.Clean()
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_history")
	defer capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")

	// Drop assistant table before test
	err := capsule.Schema().DropTableIfExists("__unit_test_conversation_assistant")
	if err != nil {
		t.Fatal(err)
	}

	conv, err := NewXun(Setting{
		Connector: "default",
		Table:     "__unit_test_conversation",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create multiple assistants for pagination testing
	for i := 0; i < 25; i++ {
		tagsJSON, err := jsoniter.MarshalToString([]string{fmt.Sprintf("tag%d", i%5)})
		if err != nil {
			t.Fatal(err)
		}

		assistant := map[string]interface{}{
			"name":        fmt.Sprintf("Assistant %d", i),
			"type":        "assistant",
			"connector":   "openai",
			"description": fmt.Sprintf("Description %d", i),
			"tags":        tagsJSON,
		}
		err = conv.SaveAssistant(assistant)
		assert.Nil(t, err)
	}

	// Test first page
	resp, err := conv.GetAssistants(AssistantFilter{
		Page:     1,
		PageSize: 10,
	})
	assert.Nil(t, err)
	assert.Equal(t, 10, len(resp.P.Items))
	assert.Equal(t, 25, resp.P.Total)
	assert.Equal(t, 3, resp.P.LastPage)

	// Test second page
	resp, err = conv.GetAssistants(AssistantFilter{
		Page:     2,
		PageSize: 10,
	})
	assert.Nil(t, err)
	assert.Equal(t, 10, len(resp.P.Items))

	// Test last page
	resp, err = conv.GetAssistants(AssistantFilter{
		Page:     3,
		PageSize: 10,
	})
	assert.Nil(t, err)
	assert.Equal(t, 5, len(resp.P.Items))

	// Test filtering with tags
	resp, err = conv.GetAssistants(AssistantFilter{
		Tags:     []string{"tag0"},
		Page:     1,
		PageSize: 10,
	})
	assert.Nil(t, err)
	assert.Equal(t, 5, len(resp.P.Items))
}
