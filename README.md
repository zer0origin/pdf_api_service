# API Backend Rework
URL Entities:
- /documents/
- /selections/
- /users/

These will be the root entities and will MUST be the root of the controllers.

- GET    200    /documents/:id	Get a document.
- POST   200 	/documents/		Create a new document.
- PUT    200 	/documents/:id	Create or update an existing document.
- DELETE 200    /documents/:id	Delete a document.

In the example above, this request is handled by the document controller as that is the main entity.

Documentation can be created by using swag init.