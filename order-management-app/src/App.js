import React, { useState } from 'react';
import OrderList from './OrderList';
import ChangeStatus from './ChangeStatus';
import Notification from './Notification';

const App = () => {
  const [orders, setOrders] = useState([
    { id: 1, status: 'Pending' },
    { id: 2, status: 'Shipped' },
    { id: 3, status: 'Delivered' },
  ]);

  const updateStatus = (orderId, newStatus) => {
    setOrders(prevOrders =>
      prevOrders.map(order =>
        order.id === orderId ? { ...order, status: newStatus } : order
      )
    );
  };

  return (
    <div>
      <OrderList orders={orders} />
      {orders.map(order => (
        <div key={order.id}>
          <ChangeStatus orderId={order.id} updateStatus={updateStatus} />
        </div>
      ))}
      <Notification message="Shipment status updated successfully" />
    </div>
  );
};

export default App;
