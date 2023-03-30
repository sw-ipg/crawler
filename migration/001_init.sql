-- +goose Up
create table if not exists public.index
(
    id                serial
    constraint index_pk_id
    primary key,
    domain            text                    not null,
    path              text                    not null,
    crc32_checksum    bigint                  not null,
    storage_file_name text                    not null,
    date              timestamp default now() not null
);

create unique index if not exists index_path_crc32_checksum_uindex
    on public.index (path, crc32_checksum);

-- +goose Down
drop table if exists public.index;

