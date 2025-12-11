import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View, TouchableOpacity, ActivityIndicator } from 'react-native';
import { useState } from 'react';
import { apiClient } from './src/lib/api';

interface ApiResponse {
  message: string;
}

interface HealthResponse {
  status: string;
}

export default function App() {
  const [message, setMessage] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const fetchWelcome = async () => {
    setLoading(true);
    try {
      const data = await apiClient.get<ApiResponse>('/');
      setMessage(data.message);
    } catch (error) {
      setMessage('Error: ' + (error as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const fetchHealth = async () => {
    setLoading(true);
    try {
      const data = await apiClient.get<HealthResponse>('/api/health');
      setMessage('Health: ' + data.status);
    } catch (error) {
      setMessage('Error: ' + (error as Error).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Mobile App</Text>
      
      <View style={styles.buttonContainer}>
        <TouchableOpacity style={styles.button} onPress={fetchWelcome}>
          <Text style={styles.buttonText}>Test API</Text>
        </TouchableOpacity>
        
        <TouchableOpacity style={styles.button} onPress={fetchHealth}>
          <Text style={styles.buttonText}>Check Health</Text>
        </TouchableOpacity>
      </View>

      {loading && <ActivityIndicator size="large" color="#2563eb" />}
      {message && <Text style={styles.message}>{message}</Text>}
      
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 30,
  },
  buttonContainer: {
    gap: 15,
    marginBottom: 20,
  },
  button: {
    backgroundColor: '#2563eb',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
    minWidth: 200,
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  message: {
    marginTop: 20,
    fontSize: 16,
    textAlign: 'center',
  },
});
