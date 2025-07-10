# API Backend Rework
API to be reworked and written in Go. API urls to be restructured to make more senses.

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