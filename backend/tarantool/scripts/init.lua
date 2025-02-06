log = require('log')  -- Добавляем логирование
clock = require('clock')

-- Конфигурируем Tarantool
box.cfg {
    listen = 3301,
    log_level = 5,
    memtx_memory = 1024 * 1024 * 1024  -- Увеличим память до 1 ГБ
}

log.info("Tarantool started!")


local username = os.getenv('TARANTOOL_USER_NAME') or 'admin'
local password = os.getenv('TARANTOOL_USER_PASSWORD') or 'Passw0rd'

if username then
    log.info("env variable TARANTOOL_USER_NAME found")
    if box.schema.user.exists(username) then
        log.info("changing password")
        box.schema.user.passwd(username, password)
    end
    if not box.schema.user.exists(username) then
        log.info("creating user")
        box.schema.user.create(username, {password = password})
        box.schema.user.grant(username, 'read,write,execute', 'universe')
    end
end

-- Создаем пространство для диалогов
box.schema.space.create('dialogs', {
    if_not_exists = true,
    format = {
        { name = 'dialog_id', type = 'string' },
        { name = 'user_id', type = 'string' }
    }
})
box.space.dialogs:create_index('primary', {
    type = 'tree',
    parts = { 'user_id', 'dialog_id' },
    if_not_exists = true
})

-- Создаем пространство для сообщений
box.schema.space.create('messages', {
    if_not_exists = true,
    format = {
        { name = 'dialog_id', type = 'string' },
        { name = 'message_id', type = 'unsigned' },
        { name = 'author_id', type = 'string' },
        { name = 'message', type = 'string' },
        { name = 'datetime', type = 'number' }
    }
})
-- 
box.space.messages:create_index('primary', {
    type = 'tree',
    parts = { 'dialog_id', 'message_id' },
    if_not_exists = true
})
-- 
box.schema.sequence.create('message_seq', { if_not_exists = true })

log.info("Spaces 'dialogs' and 'messages' created!")

-- Функция для сохранения сообщения
function save_message(my_id, partner_id, message)
    local dialog_id = my_id < partner_id and my_id .. ':' .. partner_id or partner_id .. ':' .. my_id

    -- Вставляем записи в dialogs (если записи нет, она добавится, если есть – обновится)
    box.space.dialogs:replace({dialog_id, my_id})
    box.space.dialogs:replace({dialog_id, partner_id})

    -- Вставляем сообщение
    message_id = box.sequence.message_seq:next()
    box.space.messages:insert({ dialog_id, message_id, my_id, message, clock.time() })

    -- log.info("Message saved: %s -> %s: %s", my_id, partner_id, message)

    return true
end

-- Функция для получения сообщений
function get_dialog(my_id, partner_id, limit, offset)
    local dialog_id = my_id < partner_id and my_id .. ':' .. partner_id or partner_id .. ':' .. my_id

    local result = box.space.messages.index.primary:select({dialog_id}, {
        iterator = 'REQ',
        limit = limit,
        offset = offset
    })

    -- log.info("Retrieved %d messages for dialog %s", #result, dialog_id)

    return result
end

-- Функция для получения списка диалогов пользователя
function get_dialogs(my_id, limit, offset)
    local result = box.space.dialogs.index.primary:select({my_id}, {
        iterator = 'EQ',
        limit = limit,
        offset = offset
    })

    -- log.info("Retrieved %d dialogs for user %s", #result, my_id)

    return result
end

-- Регистрируем функции
box.schema.func.create('save_message', {language = 'Lua', if_not_exists = true})
box.schema.func.create('get_dialog', {language = 'Lua', if_not_exists = true})
box.schema.func.create('get_dialogs', {language = 'Lua', if_not_exists = true})

log.info("Tarantool setup complete!")
