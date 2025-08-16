
const NotFound = () => {
  return (
    <div className="w-full min-h-full flex flex-row justify-center items-start pl-16">
      <section className="@container w-full max-w-5xl min-h-svh flex flex-col justify-start items-center pt-6 pb-8">
        <div className="w-full px-4 grow flex flex-col justify-center items-center px-6">
          <p className="font-medium">{"The page you are looking for can't be found."}</p>
          <p className="mt-4 text-[8rem] font-mono text-foreground">404</p>
        </div>
      </section>
    </div>
  );
};

export default NotFound;
