import React, { useState, useEffect } from 'react';
import axios from 'axios';

const TodoCard = props => {
  const data = props.data;
  const [itemInfo, setItemInfo] = useState();

  const postTodoStatus = () => {
    axios.post("/todo/" + data.id + "/done")
      .then(function () {
        window.location.reload();
      });
  }

  useEffect(() => {
    setItemInfo(data);
  }, [data])

  return (
    <div>
      {
        itemInfo ?
          (
            itemInfo.is_done === true ? (
              < div onClick={postTodoStatus} >☑ {itemInfo.name}</div>
            ) : (
                < div onClick={postTodoStatus} >□ {itemInfo.name}</div>
              )
          ) : (
            <div></div>
          )
      }
    </div >
  );
}

export default TodoCard;