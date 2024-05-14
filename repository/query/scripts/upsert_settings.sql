INSERT INTO settings (key, value)
VALUES (@key, @value)
ON CONFLICT (key)
DO UPDATE SET value = @value
