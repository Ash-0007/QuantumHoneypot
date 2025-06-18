import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Form, Button, Card, Alert, Spinner, Tabs, Tab, Table, Badge } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import axios from 'axios';

// Define API base URLs
const BACKEND_URL = 'http://localhost:8082';
const AI_SERVICE_URL = 'http://localhost:5000';

function App() {
  // Decoy generation state
  const [target, setTarget] = useState('');
  const [complexity, setComplexity] = useState(5);
  const [count, setCount] = useState(10);
  const [decoys, setDecoys] = useState([]);
  
  // Key generation state
  const [algorithm, setAlgorithm] = useState('kyber');
  const [decoyCount, setDecoyCount] = useState(5);
  const [keyPair, setKeyPair] = useState(null);
  
  // Encryption state
  const [message, setMessage] = useState('');
  const [publicKey, setPublicKey] = useState('');
  const [encryptedData, setEncryptedData] = useState(null);
  
  // Decryption state
  const [privateKey, setPrivateKey] = useState('');
  const [ciphertext, setCiphertext] = useState('');
  const [nonce, setNonce] = useState('');
  const [decryptedMessage, setDecryptedMessage] = useState('');
  
  // UI state
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [activeTab, setActiveTab] = useState('encryption');
  
  // Backend services status
  const [backendStatus, setBackendStatus] = useState({ go: 'unknown', ai: 'unknown' });
  
  useEffect(() => {
    // Check backend status on load
    checkStatus();
  }, []);
  
  const checkStatus = async () => {
    try {
      // Check Go backend
      const goResponse = await axios.get(`${BACKEND_URL}/api/health`);
      setBackendStatus(prev => ({ ...prev, go: 'operational' }));
      
      // Check AI service
      const aiResponse = await axios.get(`${AI_SERVICE_URL}/health`);
      setBackendStatus(prev => ({ ...prev, ai: 'operational' }));
    } catch (err) {
      console.error('Error checking status:', err);
    }
  };

  const generateDecoys = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await axios.post(`${BACKEND_URL}/api/decoys/generate`, {
        target,
        complexity: parseInt(complexity),
        count: parseInt(count)
      });
      
      setDecoys(response.data.decoys || []);
      setSuccess('Decoys generated successfully!');
    } catch (err) {
      console.error('Error generating decoys:', err);
      setError('Failed to generate decoys. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  const generateKeyPair = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await axios.post(`${BACKEND_URL}/api/keys/generate`, {
        algorithm,
        count: parseInt(decoyCount)
      });
      
      setKeyPair(response.data);
      setSuccess('Key pair generated successfully!');
      
      // Set the public key for encryption automatically
      if (response.data.public_key) {
        setPublicKey(response.data.public_key);
      }
      
      // Set the private key for decryption automatically
      if (response.data.private_key) {
        setPrivateKey(response.data.private_key);
      }
    } catch (err) {
      console.error('Error generating key pair:', err);
      setError('Failed to generate key pair. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  const encryptMessage = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await axios.post(`${BACKEND_URL}/api/encrypt`, {
        public_key: publicKey,
        algorithm,
        data: message
      });
      
      setEncryptedData(response.data);
      if (response.data.ciphertext) {
        setCiphertext(response.data.ciphertext);
      }
      if (response.data.nonce) {
        setNonce(response.data.nonce);
      }
      setSuccess('Message encrypted successfully!');
    } catch (err) {
      console.error('Error encrypting message:', err);
      setError('Failed to encrypt message. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  const decryptMessage = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await axios.post(`${BACKEND_URL}/api/decrypt`, {
        private_key: privateKey,
        ciphertext,
        nonce,
        algorithm
      });
      
      if (response.data.plaintext) {
        setDecryptedMessage(response.data.plaintext);
        setSuccess('Message decrypted successfully!');
      } else {
        setError('Decryption returned no plaintext');
      }
    } catch (err) {
      console.error('Error decrypting message:', err);
      setError('Failed to decrypt message. Please check your private key.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="py-4">
      <Row className="mb-4">
        <Col>
          <h1>Post-Quantum Cognitive Decoys</h1>
          <p className="lead">
            Enhance your cryptographic security with post-quantum algorithms and cognitive decoy generation
          </p>
          <div className="d-flex gap-2 mb-3">
            <Badge bg={backendStatus.go === 'operational' ? 'success' : 'danger'}>Go Backend: {backendStatus.go}</Badge>
            <Badge bg={backendStatus.ai === 'operational' ? 'success' : 'danger'}>AI Service: {backendStatus.ai}</Badge>
          </div>
        </Col>
      </Row>
      
      {error && <Alert variant="danger">{error}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}
      
      <Tabs
        activeKey={activeTab}
        onSelect={(k) => setActiveTab(k)}
        className="mb-3"
        fill
      >
        <Tab eventKey="encryption" title="Encryption & Decryption">
          <Row>
            <Col md={6}>
              <Card className="mb-4">
                <Card.Header>Key Generator</Card.Header>
                <Card.Body>
                  <Form onSubmit={generateKeyPair}>
                    <Form.Group className="mb-3">
                      <Form.Label>Post-Quantum Algorithm</Form.Label>
                      <Form.Select
                        value={algorithm}
                        onChange={(e) => setAlgorithm(e.target.value)}
                      >
                        <option value="kyber">Kyber (KEM)</option>
                        <option value="saber">Saber (KEM)</option>
                        <option value="ntru">NTRU (KEM)</option>
                        <option value="dilithium">Dilithium (Signature)</option>
                      </Form.Select>
                    </Form.Group>
                    
                    <Form.Group className="mb-3">
                      <Form.Label>Number of Decoy Keys</Form.Label>
                      <Form.Control
                        type="number"
                        min="0"
                        max="20"
                        value={decoyCount}
                        onChange={(e) => setDecoyCount(e.target.value)}
                      />
                    </Form.Group>
                    
                    <Button variant="primary" type="submit" disabled={loading}>
                      {loading ? (
                        <>
                          <Spinner as="span" animation="border" size="sm" className="me-2" />
                          Generating...
                        </>
                      ) : (
                        'Generate Key Pair'
                      )}
                    </Button>
                  </Form>
                  
                  {keyPair && (
                    <div className="mt-3">
                      <h5>Key Pair Details:</h5>
                      <p><strong>Algorithm:</strong> {keyPair.algorithm}</p>
                      <p><strong>Fingerprint:</strong> {keyPair.fingerprint}</p>
                      <p><strong>Decoys:</strong> {keyPair.decoys ? keyPair.decoys.length : 0}</p>
                      <p><strong>Generated:</strong> {new Date(keyPair.generated_at).toLocaleString()}</p>
                    </div>
                  )}
                </Card.Body>
              </Card>
            </Col>
            
            <Col md={6}>
              <Row>
                <Col>
                  <Card className="mb-4">
                    <Card.Header>Encryption</Card.Header>
                    <Card.Body>
                      <Form onSubmit={encryptMessage}>
                        <Form.Group className="mb-3">
                          <Form.Label>Public Key</Form.Label>
                          <Form.Control
                            as="textarea"
                            rows={3}
                            value={publicKey}
                            onChange={(e) => setPublicKey(e.target.value)}
                            placeholder="Paste public key"
                            required
                          />
                        </Form.Group>
                        
                        <Form.Group className="mb-3">
                          <Form.Label>Message</Form.Label>
                          <Form.Control
                            as="textarea"
                            rows={3}
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            placeholder="Enter message to encrypt"
                            required
                          />
                        </Form.Group>
                        
                        <Button variant="success" type="submit" disabled={loading}>
                          {loading ? (
                            <>
                              <Spinner as="span" animation="border" size="sm" className="me-2" />
                              Encrypting...
                            </>
                          ) : (
                            'Encrypt Message'
                          )}
                        </Button>
                      </Form>
                      
                      {encryptedData && (
                        <div className="mt-3">
                          <h5>Encrypted Data:</h5>
                          <p><strong>Algorithm:</strong> {encryptedData.algorithm}</p>
                          <small className="text-muted d-block mb-2">Ciphertext (first 50 chars):</small>
                          <div className="bg-light p-2 rounded mb-2 text-break">
                            {encryptedData.ciphertext ? encryptedData.ciphertext.substring(0, 50) + '...' : ''}
                          </div>
                        </div>
                      )}
                    </Card.Body>
                  </Card>
                </Col>
              </Row>
              <Row>
                <Col>
                  <Card>
                    <Card.Header>Decryption</Card.Header>
                    <Card.Body>
                      <Form onSubmit={decryptMessage}>
                        <Form.Group className="mb-3">
                          <Form.Label>Private Key</Form.Label>
                          <Form.Control
                            as="textarea"
                            rows={3}
                            value={privateKey}
                            onChange={(e) => setPrivateKey(e.target.value)}
                            placeholder="Paste private key"
                            required
                          />
                        </Form.Group>
                        
                        <Form.Group className="mb-3">
                          <Form.Label>Ciphertext</Form.Label>
                          <Form.Control
                            as="textarea"
                            rows={3}
                            value={ciphertext}
                            onChange={(e) => setCiphertext(e.target.value)}
                            placeholder="Enter ciphertext to decrypt"
                            required
                          />
                        </Form.Group>
                        
                        <Form.Group className="mb-3">
                          <Form.Label>Nonce</Form.Label>
                          <Form.Control
                            type="text"
                            value={nonce}
                            onChange={(e) => setNonce(e.target.value)}
                            placeholder="Enter nonce"
                            required
                          />
                        </Form.Group>
                        
                        <Button variant="warning" type="submit" disabled={loading}>
                          {loading ? (
                            <>
                              <Spinner as="span" animation="border" size="sm" className="me-2" />
                              Decrypting...
                            </>
                          ) : (
                            'Decrypt Message'
                          )}
                        </Button>
                      </Form>
                      
                      {decryptedMessage && (
                        <div className="mt-3">
                          <h5>Decrypted Message:</h5>
                          <div className="bg-light p-2 rounded text-break">
                            {decryptedMessage}
                          </div>
                        </div>
                      )}
                    </Card.Body>
                  </Card>
                </Col>
              </Row>
            </Col>
          </Row>
        </Tab>
        
        <Tab eventKey="decoys" title="Decoy Generation">
          <Row>
            <Col md={6}>
              <Card className="mb-4">
                <Card.Header>Decoy Generator</Card.Header>
                <Card.Body>
                  <Form onSubmit={generateDecoys}>
                    <Form.Group className="mb-3">
                      <Form.Label>Target Text</Form.Label>
                      <Form.Control
                        type="text"
                        placeholder="Enter target text (e.g., kyber768)"
                        value={target}
                        onChange={(e) => setTarget(e.target.value)}
                        required
                      />
                      <Form.Text className="text-muted">
                        This is the text for which decoys will be generated
                      </Form.Text>
                    </Form.Group>
                    
                    <Form.Group className="mb-3">
                      <Form.Label>Complexity ({complexity})</Form.Label>
                      <Form.Range
                        min="1"
                        max="10"
                        value={complexity}
                        onChange={(e) => setComplexity(e.target.value)}
                      />
                      <Form.Text className="text-muted">
                        Higher complexity = more different from target
                      </Form.Text>
                    </Form.Group>
                    
                    <Form.Group className="mb-3">
                      <Form.Label>Number of Decoys</Form.Label>
                      <Form.Control
                        type="number"
                        min="1"
                        max="100"
                        value={count}
                        onChange={(e) => setCount(e.target.value)}
                      />
                    </Form.Group>
                    
                    <Button variant="primary" type="submit" disabled={loading}>
                      {loading ? (
                        <>
                          <Spinner as="span" animation="border" size="sm" className="me-2" />
                          Generating...
                        </>
                      ) : (
                        'Generate Decoys'
                      )}
                    </Button>
                  </Form>
                </Card.Body>
              </Card>
            </Col>
            
            <Col md={6}>
              <Card>
                <Card.Header>Generated Decoys</Card.Header>
                <Card.Body>
                  {decoys.length > 0 ? (
                    <Table striped hover>
                      <thead>
                        <tr>
                          <th>#</th>
                          <th>Decoy</th>
                          <th>Similarity</th>
                        </tr>
                      </thead>
                      <tbody>
                        {decoys.map((decoy, index) => (
                          <tr key={index}>
                            <td>{index + 1}</td>
                            <td>{typeof decoy === 'string' ? decoy : decoy.value || 'N/A'}</td>
                            <td>{typeof decoy === 'object' && decoy.similarity ? `${(decoy.similarity * 100).toFixed(1)}%` : 'N/A'}</td>
                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  ) : (
                    <p className="text-muted">No decoys generated yet</p>
                  )}
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </Tab>
        
        <Tab eventKey="about" title="About">
          <Row>
            <Col>
              <Card>
                <Card.Header>About Post-Quantum Cognitive Decoys</Card.Header>
                <Card.Body>
                  <h4>What is Post-Quantum Cryptography?</h4>
                  <p>
                    Post-quantum cryptography refers to cryptographic algorithms that are thought to be secure against an attack by a quantum computer. As quantum computing advances, it threatens many commonly used cryptographic algorithms.
                  </p>
                  
                  <h4>What are Cognitive Decoys?</h4>
                  <p>
                    Cognitive decoys are false artifacts designed to appear genuine to attackers. In cryptography, cognitive decoys can be false keys, certificates, or other cryptographic material that waste an attacker's resources.
                  </p>
                  
                  <h4>Implementation Details</h4>
                  <p>
                    This application implements a functional simulation of post-quantum cryptography algorithms:
                  </p>
                  <ul>
                    <li><strong>Kyber:</strong> Lattice-based key encapsulation mechanism (KEM)</li>
                    <li><strong>Saber:</strong> Module-LWR based KEM</li>
                    <li><strong>NTRU:</strong> Lattice-based cryptosystem</li>
                    <li><strong>Dilithium:</strong> Lattice-based digital signature algorithm</li>
                  </ul>
                  
                  <div className="alert alert-info">
                    <strong>Note:</strong> This is a demonstration implementation. For production use, integrate with standard post-quantum libraries like CloudFlare's CIRCL, Open Quantum Safe (liboqs), or libraries containing NIST PQC standardized algorithms.
                  </div>
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </Tab>
      </Tabs>
      
      <footer className="mt-5 text-center text-muted">
        <small>Post-Quantum Cognitive Decoys (PQCD) - A demonstration of post-quantum cryptography with cognitive decoy generation</small>
      </footer>
    </Container>
  );
}

export default App; 