import React, { useState, useEffect } from 'react';
import axios from "axios";
import './App.css';
import TodoCard from "./Atoms/ItemCard";
import TodoForm from "./Atoms/Form";

axios.defaults.headers.post['Access-Control-Allow-Origin'] = '*';

const getTodos = (endpoint) => {
  if (process.env.NODE_ENV === "production") {
    endpoint = process.env.REACT_APP_API_ENDPOINT + endpoint;
  }
  return axios.get(endpoint);
}

const App = () => {
  const [todos, setTodos] = useState();

  useEffect(() => {
    getTodos("/todo")
      .then(function (response) {
        setTodos(response.data);
      })
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <div>
          <TodoForm></TodoForm>
        </div>
        <div>
          {todos ?
            (todos.length !== 0 ?
              (
                todos.map(todo => (
                  <TodoCard data={todo} />
                ))
              ) : (<div>Not found</div>)
            ) : (<div>Not found</div>)
          }
        </div>
      </header>
    </div>
  );
}

export default App;
