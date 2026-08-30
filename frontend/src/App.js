import React, { useState, useEffect } from 'react';

// url del backend. en la fase 2 esto cambia y apunta al middleware
const API_URL = 'http://localhost:8080';

function App() {
  const [nodos, setNodos] = useState([]);
  const [nombre, setNombre] = useState('');
  const [tipo, setTipo] = useState('Router');
  const [ip, setIp] = useState('');

  const [origenEnlace, setOrigenEnlace] = useState('');
  const [destinoEnlace, setDestinoEnlace] = useState('');

  const [origenPaquete, setOrigenPaquete] = useState('');
  const [destinoPaquete, setDestinoPaquete] = useState('');
  const [resultado, setResultado] = useState(null);

  // trae los nodos apenas carga la pagina
  useEffect(() => {
    cargarNodos();
  }, []);

  function cargarNodos() {
    fetch(API_URL + '/nodos')
      .then(res => res.json())
      .then(data => setNodos(data || []));
  }

  function crearNodo(e) {
    e.preventDefault();
    fetch(API_URL + '/nodos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nombre, tipo, ip })
    })
      .then(res => res.json())
      .then(() => {
        setNombre('');
        setIp('');
        cargarNodos(); // recargamos todo, no es lo mas eficiente pero jala
      });
  }

  function crearEnlace(e) {
    e.preventDefault();
    fetch(API_URL + '/enlaces', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        nodo_origen_id: parseInt(origenEnlace),
        nodo_destino_id: parseInt(destinoEnlace)
      })
    }).then(() => alert('enlace creado!'));
  }

  function enviarPaquete(e) {
    e.preventDefault();
    setResultado(null);
    fetch(API_URL + '/enviar', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        origen: parseInt(origenPaquete),
        destino: parseInt(destinoPaquete)
      })
    })
      .then(res => res.json())
      .then(data => setResultado(data))
      .catch(err => setResultado({ error: 'algo tronó: ' + err }));
  }

  return (
    <div style={{ fontFamily: 'Arial', margin: '20px' }}>
      <h1>Simulador de Red (Proyecto Cómputo Distribuido) - Fase 1</h1>

      <h2>1. Crear Nodo</h2>
      <form onSubmit={crearNodo}>
        <input placeholder="nombre" value={nombre} onChange={e => setNombre(e.target.value)} required />
        <select value={tipo} onChange={e => setTipo(e.target.value)}>
          <option value="Router">Router</option>
          <option value="Switch">Switch</option>
          <option value="Server">Server</option>
          <option value="Endpoint">Endpoint</option>
        </select>
        <input placeholder="ip (opcional)" value={ip} onChange={e => setIp(e.target.value)} />
        <button type="submit">Agregar Nodo</button>
      </form>

      <h2>2. Nodos registrados</h2>
      <ul>
        {nodos.map(n => (
          <li key={n.id}>#{n.id} - {n.nombre} ({n.tipo}) {n.ip}</li>
        ))}
      </ul>

      <h2>3. Conectar Nodos</h2>
      <form onSubmit={crearEnlace}>
        <select value={origenEnlace} onChange={e => setOrigenEnlace(e.target.value)} required>
          <option value="">-- origen --</option>
          {nodos.map(n => <option key={n.id} value={n.id}>{n.nombre}</option>)}
        </select>
        <select value={destinoEnlace} onChange={e => setDestinoEnlace(e.target.value)} required>
          <option value="">-- destino --</option>
          {nodos.map(n => <option key={n.id} value={n.id}>{n.nombre}</option>)}
        </select>
        <button type="submit">Conectar</button>
      </form>

      <h2>4. Enviar Paquete</h2>
      <form onSubmit={enviarPaquete}>
        <select value={origenPaquete} onChange={e => setOrigenPaquete(e.target.value)} required>
          <option value="">-- origen --</option>
          {nodos.map(n => <option key={n.id} value={n.id}>{n.nombre}</option>)}
        </select>
        <select value={destinoPaquete} onChange={e => setDestinoPaquete(e.target.value)} required>
          <option value="">-- destino --</option>
          {nodos.map(n => <option key={n.id} value={n.id}>{n.nombre}</option>)}
        </select>
        <button type="submit">Simular Envío</button>
      </form>

      {resultado && (
        <div style={{ marginTop: '10px', border: '1px solid black', padding: '10px' }}>
          {resultado.error ? (
            <p style={{ color: 'red' }}>Error: {resultado.error}</p>
          ) : (
            <p>Ruta tomada: {resultado.ruta}</p>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
