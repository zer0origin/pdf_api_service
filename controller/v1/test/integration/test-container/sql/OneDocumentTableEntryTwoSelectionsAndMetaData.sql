insert into document_table ("Document_UUID", "Document_Base64", "Owner_UUID", "Owner_Type")
values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 'Fake document for testing', uuid('f701aa7e-10e9-48b9-83f1-6b035a5b7564'), 1);

insert into selection_table ("Selection_UUID", "Document_UUID")
values (uuid('a5fdea38-0a86-4c19-ae4f-c87a01bc860d'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'));

insert into selection_table ("Selection_UUID", "Document_UUID")
values (uuid('335a6b95-6707-4e2b-9c37-c76d017f6f97'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'));

insert into documentmeta_table ("Document_UUID", "Number_Of_Pages", "Height", "Width", "Images")
values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 31, 1920, 1080, '{}');
