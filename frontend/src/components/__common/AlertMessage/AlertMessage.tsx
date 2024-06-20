import React from 'react';
import { AlertMessageType } from '@hooks/useAlertMessage';
import { InformationCircleIcon } from '@heroicons/react/24/solid';
interface AlertMessageProps {
  message: string;
  type: AlertMessageType;
  onClose: () => void;
}

const AlertMessage: React.FC<AlertMessageProps> = ({ message, type, onClose }) => {
  const getAlertClasses = (type: AlertMessageType) => {
    switch (type) {
      case AlertMessageType.SUCCESS:
        return 'bg-green-500 text-white';
      case AlertMessageType.WARNING:
        return 'bg-yellow-500 text-white';
      case AlertMessageType.ERROR:
        return 'bg-red-500 text-white';
      case AlertMessageType.INFO:
      default:
        return 'bg-blue-500 text-white';
    }
  };

  return (
    <div className={`fixed top-0 left-1/3 right-1/3 z-50 p-4 ${getAlertClasses(type)}`}>
      <div className="container mx-auto w-72 flex items-center justify-center px-6 py-4">
        <div className="flex items-center">
          <InformationCircleIcon width={'10%'}/>
          <span className="alert-message font-semibold ml-2">{message}</span>
        </div>
        <button className="ml-4 text-lg font-semibold outline-none focus:outline-none" onClick={onClose}>
          <span className="text-white">×</span>
        </button>
      </div>
    </div>
  );
};

export default AlertMessage;
