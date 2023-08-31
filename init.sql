-- сгенерировано с GORM AutoMigrate
create table users
(
    id bigserial not null
        constraint users_pkey
            primary key,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create table segments
(
    id bigserial not null
        constraint segments_pkey
            primary key,
    slug text,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index idx_segments_slug
    on segments (slug);

create table user_segments
(
    user_id bigint
        constraint fk_user_segments_user
            references users,
    segment_id bigint
        constraint fk_user_segments_segment
            references segments,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone
);

