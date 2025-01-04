-- Description: Create a trigger to notify when a new entry is inserted

CREATE FUNCTION notify_insert()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_entry', NEW.deviceuuid);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notify_insert
AFTER INSERT ON entries
FOR EACH ROW EXECUTE FUNCTION notify_insert();