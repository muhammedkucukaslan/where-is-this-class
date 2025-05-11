CREATE TABLE class_rooms (
    id UUID PRIMARY KEY,
    code TEXT UNIQUE,
    floor INT,
    image_url TEXT,
    visited INT DEFAULT 0
);


CREATE TABLE class_room_translations (
    id UUID PRIMARY KEY,
    class_room_id UUID REFERENCES class_rooms(id) ON DELETE CASCADE,
    language TEXT,
    building TEXT,
    description TEXT,
    detail Text ,
    CONSTRAINT unique_classroom_language UNIQUE (id, language)  
);
