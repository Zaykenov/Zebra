import { Column } from "react-table";

export default function extractHeadersAndAccessor(columns: any[]) {
    const headers: any[] = [];
    const accessors: any[] = [];
  
    columns.forEach(column => {
      headers.push(column.Header);
      accessors.push(column.accessor);
    });
  
    return { headers, accessors };
}