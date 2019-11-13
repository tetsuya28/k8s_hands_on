import React, { useState } from 'react';
import { FormControl, InputLabel, Input, FormHelperText, Button } from '@material-ui/core';
import makeStyles from "@material-ui/core/styles/makeStyles";
import axios from 'axios';

const useStyle = makeStyles({
  form: {
    color: "white"
  },
  button: {
    margin: 10
  }
})

const TodoForm = props => {
  const classes = useStyle(props)
  const [todoName, setTodoName] = useState("");

  const changeTodo = e => {
    setTodoName(e.target.value);
  }

  const postTodo = () => {
    if (todoName === "") {
      alert("未入力です");
      return false;
    }
    const params = {
      name: todoName
    }
    axios.post("/todo", params)
      .then(function () {
        window.location.reload();
      });
  }

  return (
    <div>
      <FormControl>
        <InputLabel htmlFor="todo_name" className={classes.form}>Task Name</InputLabel>
        <Input id="todo_name" aria-describedby="my-helper-text" onChange={changeTodo} className={classes.form} required />
        <FormHelperText id="my-helper-text" className={classes.form}>タスクの名前を入力してください</FormHelperText>
        <Button variant="outlined" color="inherit" onClick={postTodo} className={classes.button}>Submit</Button>
      </FormControl>
    </div>
  );
}

export default TodoForm;