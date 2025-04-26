
--- create message 
CREATE OR REPLACE FUNCTION create_message(_user_id BIGINT, _content TEXT)
RETURNS UUID AS $$
DECLARE
    _id UUID;
BEGIN
    INSERT INTO messages (user_id, content)
    VALUES (_user_id, _content)
    RETURNING id INTO _id;

    RETURN _id;
END;
$$ LANGUAGE plpgsql;

--- get all messages
CREATE OR REPLACE FUNCTION get_user_messages(_user_id BIGINT)
RETURNS TABLE (
    id UUID,
    content TEXT,
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT id, content, created_at
    FROM messages
    WHERE user_id = _user_id
    ORDER BY created_at DESC;
END;
$$ LANGUAGE plpgsql;

--- delete message by id
CREATE OR REPLACE FUNCTION delete_message(_user_id BIGINT, _message_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    DELETE FROM messages
    WHERE id = _message_id AND user_id = _user_id;

    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- delete account
CREATE OR REPLACE FUNCTION delete_account(_user_id BIGINT)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users WHERE id = _user_id;
END;
$$ LANGUAGE plpgsql;

--- export user data
CREATE OR REPLACE FUNCTION export_user_data(_user_id BIGINT)
RETURNS JSONB AS $$
DECLARE
    result JSONB;
BEGIN
    SELECT jsonb_build_object(
        'email', u.email,
        'username', u.username,
        'created_at', u.created_at,
        'messages', COALESCE((
            SELECT jsonb_agg(jsonb_build_object(
                'id', m.id,
                'content', m.content,
                'created_at', m.created_at
            ))
            FROM messages m WHERE m.user_id = u.id
        ), '[]'::jsonb) -- 👈 Force [] if no messages
    ) INTO result
    FROM users u
    WHERE u.id = _user_id;

    RETURN result;
END;
$$ LANGUAGE plpgsql;



---- auth ---

--- create user
CREATE OR REPLACE FUNCTION create_user(_username TEXT, _email TEXT, _password_hash TEXT)
RETURNS BIGINT AS $$
DECLARE
    _id BIGINT;
BEGIN
    INSERT INTO users (username, email, password_hash)
    VALUES (_username, _email, _password_hash)
    RETURNING id INTO _id;

    RETURN _id;
END;
$$ LANGUAGE plpgsql;

--- get user 
CREATE OR REPLACE FUNCTION get_user_by_email(_email TEXT)
RETURNS TABLE (
    id BIGINT,
    password_hash TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT users.id, users.password_hash
    FROM users
    WHERE users.email = _email;
END;
$$ LANGUAGE plpgsql;

--- get user by id for injection to context
CREATE OR REPLACE FUNCTION get_user_by_id(_id BIGINT)
RETURNS TABLE(
    id BIGINT,
    username TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT users.id, users.username
    FROM users
    WHERE users.id = _id;
END;
$$ LANGUAGE plpgsql;

--- delete user
CREATE OR REPLACE FUNCTION delete_user(_user_id BIGINT)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users WHERE id = _user_id;
END;
$$ LANGUAGE plpgsql;

--- delete user messages
CREATE OR REPLACE FUNCTION delete_user_messages(_user_id BIGINT)
RETURNS VOID AS $$
BEGIN
    DELETE FROM messages
    WHERE user_id = _user_id;
END;
$$ LANGUAGE plpgsql; 