import React, { useState, useEffect } from 'react';

const TodoCard = props => {
  const data = props.data;
  const [itemInfo, setItemInfo] = useState();

  const showTodoInfo = () => {
    console.log(itemInfo);
  }

  useEffect(() => {
    setItemInfo(data);
  }, [data])

  return (
    <div>
      {
        itemInfo ?
          (<div onClick={showTodoInfo} > {itemInfo.name}</div>) : (<div></div>)
      }
    </div>
  );
}

export default TodoCard;