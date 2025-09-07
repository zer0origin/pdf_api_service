insert into document_table ("Document_UUID", "Document_Base64", "Owner_UUID", "Owner_Type")
values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 'Fake document for testing', uuid('f701aa7e-10e9-48b9-83f1-6b035a5b7564'), 1);

insert into documentmeta_table ("Document_UUID", "Number_Of_Pages", "Height", "Width", "Images")
values (uuid('b66fd223-515f-4503-80cc-2bdaa50ef474'), 31, 1920, 1080, '{"0":"test0","1":"test1","2":"test2","3":"test3"}');
