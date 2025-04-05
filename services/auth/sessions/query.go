package sessions

var (
	QueryInsert = `
	INSERT INTO sessions 
		(
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING 
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
	`

	QueryCheckValidity = `
	SELECT
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
	FROM sessions
	WHERE session_id = $1 AND user_id = $2
	`

	QueryGetBySessionID = `
	SELECT
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
	FROM sessions
	WHERE session_id = $1
	`

	QueryGetByUserID = `
	SELECT
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
	FROM sessions
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	QueryInvalidateSession = `
	DELETE FROM sessions
	WHERE session_id = $1
	`

	QueryInvalidateAllUserSessions = `
	DELETE FROM sessions
	WHERE user_id = $1
	`

	QueryUpdateLastActivity = `
	UPDATE sessions
	SET last_activity = $1, updated_at = $1
	WHERE session_id = $2
	`

	QueryGetByAccessToken = `
	SELECT
		session_id, 
		user_id,
		access_token, 
		ip_address,
		user_agent,
		device_info,
		last_activity,
		is_access,
		expires_at, 
		created_at, 
		updated_at
	FROM sessions
	WHERE access_token = $1
	`
)
