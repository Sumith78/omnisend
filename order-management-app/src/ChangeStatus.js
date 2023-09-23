import React, { useState } from 'react';

const ChangeStatus = ({ orderId, updateStatus }) => {
  const [newStatus, setNewStatus] = useState('');

  const handleChange = e => {
    setNewStatus(e.target.value);
  };

  const handleSubmit = e => {
    e.preventDefault();
    updateStatus(orderId, newStatus);
    setNewStatus('');
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Change Status:
        <input type="text" value={newStatus} onChange={handleChange} />
      </label>
      <button type="submit">Submit</button>
    </form>
  );
};

export default ChangeStatus;
