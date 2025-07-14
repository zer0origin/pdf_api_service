create table if not exists document_table
(
    "Document_UUID"   uuid not null
        constraint "Document_Table_pk"
            primary key,
    "Document_Base64" text not null
);

create table if not exists documentmeta_table
(
    "Document_UUID"   uuid not null
        constraint documentmeta_table_pk
            primary key
        constraint documentmeta_table_document_table_null_fk
            references document_table
            on delete cascade,
    "Number_Of_Pages" integer,
    "Height"          numeric,
    "Width"           numeric,
    "Images"          json
);

create table if not exists selection_table
(
    "Selection_UUID"   uuid not null
        constraint selection_table_pk
            primary key,
    "Document_UUID"    uuid
        constraint "selection_table_document_table_Document_UUID_fk"
            references document_table
            on delete cascade,
    "isCompleted"      boolean default false,
    "Settings"         json,
    "Selection_bounds" json,
    "Page_Words"       json
);

insert into document_table values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 'Fake document for testing');
insert into selection_table values (uuid('a5fdea38-0a86-4c19-ae4f-c87a01bc860d'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), false, '{}', '{}', '{}');
insert into selection_table values (uuid('335a6b95-6707-4e2b-9c37-c76d017f6f97'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), false, '{}', '{}', '{}');
