-- tablas de la base de datos, idelamante crecerán.

CREATE TABLE nodos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    tipo VARCHAR(50) NOT NULL, -- Router, Switch, Server, Endpoint
    ip VARCHAR(50)
);

CREATE TABLE enlaces (
    id SERIAL PRIMARY KEY,
    nodo_origen_id INTEGER NOT NULL,
    nodo_destino_id INTEGER NOT NULL
);

CREATE TABLE logs_paquetes (
    id SERIAL PRIMARY KEY,
    origen INTEGER NOT NULL,
    destino INTEGER NOT NULL,
    ruta TEXT,
    fecha TIMESTAMP DEFAULT NOW()
);
