/*
     Migration script
 */

create table if not exists public.short_uris (
    id varchar(50) not null,
    original_url varchar(4096) not null,
    key varchar(50) not null,
    create_user varchar(50) null,
    created timestamptz null,
    update_user varchar(50) null,
    updated timestamptz null,
    constraint short_uris_pk primary key(id),
    constraint short_uris_uk unique(key)
)
;
