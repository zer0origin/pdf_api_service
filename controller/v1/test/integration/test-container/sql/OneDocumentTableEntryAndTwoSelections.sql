insert into document_table ("Document_UUID", "Document_Base64") values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 'Fake document for testing');
insert into selection_table ("Selection_UUID", "Document_UUID", "isCompleted") values (uuid('a5fdea38-0a86-4c19-ae4f-c87a01bc860d'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), false);
insert into selection_table ("Selection_UUID", "Document_UUID", "isCompleted") values (uuid('335a6b95-6707-4e2b-9c37-c76d017f6f97'), uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), false);
