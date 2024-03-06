import axios from 'axios';

class FakeServer {
  constructor() {
    this.axiosInstance = axios.create({
      baseURL: 'http://localhost:3004', // Replace with your json-server URL
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  async getRequest(url) {
    try {
      const response = await this.axiosInstance.get(url);
      return response.data;
    } catch (error) {
      throw new Error(`GET request failed: ${error.message}`);
    }
  }

  async postRequest(url, data) {
    try {
      const response = await this.axiosInstance.post(url, data);
      return response.data;
    } catch (error) {
      throw new Error(`POST request failed: ${error.message}`);
    }
  }

  async putRequest(url, data) {
    try {
      const response = await this.axiosInstance.put(url, data);
      return response.data;
    } catch (error) {
      throw new Error(`PUT request failed: ${error.message}`);
    }
  }

  async deleteRequest(url) {
    try {
      const response = await this.axiosInstance.delete(url);
      return response.data;
    } catch (error) {
      throw new Error(`DELETE request failed: ${error.message}`);
    }
  }

  setBaseURL(baseURL) {
    this.axiosInstance.defaults.baseURL = baseURL;
  }
}

const fakeServer = new FakeServer();

export default fakeServer;
