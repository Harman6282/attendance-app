CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_name VARCHAR(255) NOT NULL,
    teacher_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_teacher
        FOREIGN KEY (teacher_id)
        REFERENCES users(id)
);

CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL,
    student_id UUID NOT NULL,
    status VARCHAR(10) NOT NULL,

    CONSTRAINT fk_class
        FOREIGN KEY (class_id)
        REFERENCES classes(id),

    CONSTRAINT fk_student
        FOREIGN KEY (student_id)
        REFERENCES users(id)
);

