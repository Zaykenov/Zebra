import React, { FC, ReactNode } from "react";
import TerminalHeader from "../../__modules/TerminalHeader";

export interface TerminalLayout {
  children: ReactNode;
}

const TerminalLayout: FC<TerminalLayout> = ({ children }) => {
  return (
    <div className="h-screen flex flex-col">
      <TerminalHeader />
      <main className="grow flex overflow-hidden">{children}</main>
    </div>
  );
};

export default TerminalLayout;
