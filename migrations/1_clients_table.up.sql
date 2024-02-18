CREATE TABLE IF NOT EXISTS clientes (
    id_cliente SERIAL PRIMARY KEY,
    limite DECIMAL(10, 2) NOT NULL,
    saldo DECIMAL(10, 2) NOT NULL DEFAULT 0
);

INSERT INTO clientes (id_cliente, limite, saldo) VALUES
(1, 100000, 0),
(2, 80000, 0),
(3, 1000000, 0),
(4, 10000000, 0),
(5, 500000, 0)
ON CONFLICT (id_cliente) DO NOTHING;
