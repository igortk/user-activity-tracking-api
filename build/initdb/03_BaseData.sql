INSERT INTO events (user_id, event_action_timestamp, action, metadata, created_at)
VALUES
    (1, '2025-10-10 21:30:00+03', 'created', '{"ip": "192.168.1.10", "device": "iPhone 11", "location": "Kyiv, Ukraine"}', '2025-10-11 15:43:10+03'),
    (1, '2025-10-10 21:30:00+03', 'updated', '{"ip": "192.168.1.10", "device": "iPhone 12", "location": "Lviv, Ukraine"}', '2025-10-13 15:43:10+03'),
    (1, '2025-10-10 21:30:00+03', 'deleted', '{"ip": "192.168.1.10", "device": "iPhone 13", "location": "Kyiv, Ukraine"}', '2025-10-10 15:43:10+03'),
    (1, '2025-10-10 21:30:00+03', 'viewed', '{"ip": "192.168.1.10", "device": "iPhone 14", "location": "Kyiv, Ukraine"}', '2025-10-11 15:43:10+03');
