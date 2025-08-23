
export interface LibroYearMonthReport {
  id: number;
  yearMonth: string;
  count: number;
}

export interface Exchange {
  id?: number;
  from: string;
  to: string;
  value: string;
  createdAt?: Date;
  updatedAt?: Date;
}
